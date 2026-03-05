import { useState } from 'react';
import { cn } from '../lib/utils';
import { ServiceCard } from './ServiceCard';
import type { GroupWithStatus } from '../types';

interface ServiceGroupProps {
  group: GroupWithStatus;
}

export function ServiceGroup({ group }: ServiceGroupProps) {
  const [collapsed, setCollapsed] = useState(group.collapsed ?? false);

  const layoutClasses = {
    grid: 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5',
    list: 'flex flex-col gap-4',
    compact: 'grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4',
  };

  return (
    <section className="mb-10">
      {/* Group Header */}
      <header className="flex items-center justify-between mb-5">
        <div className="flex items-center gap-4">
          {/* Pure white text - no gradient colors */}
          <h2 className="text-2xl font-bold text-white/90">
            {group.name}
          </h2>
          <span className="px-3 py-1 rounded-full bg-white/5 border border-white/10 text-xs font-medium text-white/40 backdrop-blur-sm">
            {group.services.length}
          </span>
        </div>

        <button
          onClick={() => setCollapsed(!collapsed)}
          className="p-2.5 rounded-xl glass-button w-10 h-10 flex items-center justify-center"
          aria-label={collapsed ? 'Expand' : 'Collapse'}
          aria-expanded={!collapsed}
        >
          <svg
            className={cn(
              'w-5 h-5 text-white/40 transition-transform duration-300',
              collapsed ? 'rotate-180' : ''
            )}
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={1.5}
              d="M19 9l-7 7-7-7"
            />
          </svg>
        </button>
      </header>

      {/* Services */}
      <div
        className={cn(
          'transition-all duration-500 overflow-hidden',
          collapsed ? 'max-h-0 opacity-0' : 'max-h-full opacity-100'
        )}
      >
        <div className={cn(layoutClasses[group.layout as keyof typeof layoutClasses] || layoutClasses.grid)}>
          {group.services.map((service, idx) => (
            <div
              key={idx}
              className="animate-in fade-in slide-in-from-bottom-4"
              style={{ animationDelay: `${idx * 50}ms` }}
            >
              <ServiceCard service={service} />
            </div>
          ))}
        </div>
      </div>

      {/* Empty state */}
      {!collapsed && group.services.length === 0 && (
        <div className="text-center py-12 glass-card">
          <p className="text-white/40">No services in this group</p>
        </div>
      )}
    </section>
  );
}
