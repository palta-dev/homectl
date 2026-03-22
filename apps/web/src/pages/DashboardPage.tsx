import { useState, useMemo } from 'react';
import { TopBar } from '../components/TopBar';
import { Dashboard } from '../components/Dashboard';
import { useServices } from '../hooks/api';
import { useSearchStore } from '../stores/searchStore';

export function DashboardPage() {
  const { setQuery } = useSearchStore();
  const { data: services } = useServices();
  const [searchQuery, setSearchQuery] = useState('');

  if (services) {
    console.log('API Services Response:', services);
  }

  const filteredGroups = useMemo(() => {
    if (!services?.groups) return [];
    if (!searchQuery.trim()) return services.groups;

    const searchLower = searchQuery.toLowerCase().trim();
    
    return services.groups
      .map((group) => ({
        ...group,
        services: group.services.filter(
          (service) =>
            service.name.toLowerCase().includes(searchLower) ||
            service.description?.toLowerCase().includes(searchLower) ||
            service.tags?.some((tag) => tag.toLowerCase().includes(searchLower))
        ),
      }))
      .filter((group) => group.services.length > 0);
  }, [services, searchQuery]);

  const systemWidgets = useMemo(() => {
    const group = services?.groups?.find(g => g.name === 'System Health');
    return group?.services[0]?.widgets || [];
  }, [services]);

  return (
    <>
      <TopBar 
        searchQuery={searchQuery}
        systemWidgets={systemWidgets}
        onSearchChange={(q: string) => {
          setSearchQuery(q);
          setQuery(q);
        }}
      />
      <Dashboard 
        groups={filteredGroups}
        searchQuery={searchQuery}
      />
    </>
  );
}
