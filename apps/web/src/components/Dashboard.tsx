import { ServiceCard } from './ServiceCard';
import type { GroupWithStatus } from '../types';

interface DashboardProps {
  groups: GroupWithStatus[];
  searchQuery: string;
}

export function Dashboard({ groups, searchQuery }: DashboardProps) {
  // Extract system health group if it exists
  const mainGroups = groups.filter(g => g.name !== 'System Health');

  return (
    <div className="dashboard">
      {/* Service Groups */}
      {mainGroups.length > 0 ? (
        mainGroups.map((group, idx) => (
          <section key={idx} className="section">
            <div className="section-head">
              <h2>{group.name}</h2>
              <span className="badge">{group.services.length} services</span>
              <div className="divider"></div>
            </div>
            <div className="grid">
              {group.services.map((service, sIdx) => (
                <ServiceCard key={sIdx} service={service} />
              ))}
            </div>
          </section>
        ))
      ) : (
        <div className="empty-state">
          <div className="empty-icon">🔍</div>
          <h2>No services found</h2>
          <p>
            {searchQuery 
              ? `No services match "${searchQuery}". Try a different search term.`
              : 'Add services to your configuration to get started.'}
          </p>
        </div>
      )}

<footer className="dashboard-footer">
  <p>
    <a
      href="https://github.com/0bfu5c4t3/homectl"
      target="_blank"
      rel="noopener noreferrer"
    >
      homectl
    </a>{" "}
    — a product of{" "}
    <a
      href="https://palta.cloud"
      target="_blank"
      rel="noopener noreferrer"
    >
      Palta
    </a>
  </p>
</footer>
    </div>
  );
}
