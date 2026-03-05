import { TopBar } from '../components/TopBar';
import { useSearchStore } from '../stores/searchStore';
import { useState } from 'react';

// Mock incident data - will be connected to API later
const mockIncidents = [
  {
    id: 1,
    service: 'Grafana',
    severity: 'down',
    started: '2024-03-02T10:30:00Z',
    ended: '2024-03-02T11:45:00Z',
    duration: '1h 15m',
    error: 'Connection refused: port 3000',
  },
  {
    id: 2,
    service: 'Prometheus',
    severity: 'degraded',
    started: '2024-03-02T14:00:00Z',
    ended: null,
    duration: 'Ongoing',
    error: 'High latency detected (>500ms)',
  },
  {
    id: 3,
    service: 'Portainer',
    severity: 'down',
    started: '2024-03-01T08:15:00Z',
    ended: '2024-03-01T08:30:00Z',
    duration: '15m',
    error: 'Service unavailable',
  },
];

export function IncidentsPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const { setQuery } = useSearchStore();

  const filteredIncidents = mockIncidents.filter(inc => 
    inc.service.toLowerCase().includes(searchQuery.toLowerCase()) ||
    inc.error.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const stats = {
    total: mockIncidents.length,
    ongoing: mockIncidents.filter(i => !i.ended).length,
    resolved: mockIncidents.filter(i => i.ended).length,
    avgDuration: '45m',
  };

  return (
    <>
      <TopBar 
        searchQuery={searchQuery}
        onSearchChange={(q) => {
          setSearchQuery(q);
          setQuery(q);
        }}
      />

      <div className="dashboard">
        {/* Stats Strip */}
        <div className="stats">
          <div className="stat">
            <span>Total</span>
            <b>{stats.total}</b>
          </div>
          <div className="stat">
            <span>Ongoing</span>
            <b style={{ color: 'var(--red)' }}>{stats.ongoing}</b>
          </div>
          <div className="stat">
            <span>Resolved</span>
            <b style={{ color: 'var(--green)' }}>{stats.resolved}</b>
          </div>
          <div className="stat">
            <span>Avg Duration</span>
            <b>{stats.avgDuration}</b>
          </div>
        </div>

        {/* Incidents Table */}
        <section className="section">
          <div className="section-head">
            <h2>Incident History</h2>
            <div className="divider"></div>
          </div>

          <table className="table">
            <thead>
              <tr>
                <th>Service</th>
                <th>Severity</th>
                <th>Started</th>
                <th>Ended</th>
                <th>Duration</th>
                <th>Error</th>
              </tr>
            </thead>
            <tbody>
              {filteredIncidents.map((inc) => (
                <tr key={inc.id}>
                  <td>
                    <b style={{ color: 'var(--text)' }}>{inc.service}</b>
                  </td>
                  <td>
                    <span className={`sev ${inc.severity}`}>
                      <span className={`b-dot ${inc.severity === 'down' ? 'off' : inc.severity === 'degraded' ? 'deg' : 'ok'}`}></span>
                      {inc.severity === 'down' ? 'Offline' : inc.severity === 'degraded' ? 'Degraded' : 'Online'}
                    </span>
                  </td>
                  <td>{new Date(inc.started).toLocaleString()}</td>
                  <td>{inc.ended ? new Date(inc.ended).toLocaleString() : '—'}</td>
                  <td>{inc.duration}</td>
                  <td>
                    <details className="errdetails">
                      <summary>
                        <span style={{ fontSize: '11px', opacity: 0.7 }}>View error</span>
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                          <polyline points="6 9 12 15 18 9"/>
                        </svg>
                      </summary>
                      <pre className="code">{inc.error}</pre>
                    </details>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          {filteredIncidents.length === 0 && (
            <div className="empty-state">
              <div className="empty-icon">📋</div>
              <h2>No incidents found</h2>
              <p>{searchQuery ? 'Try a different search term.' : 'All services are running smoothly!'}</p>
            </div>
          )}
        </section>

        <footer className="dashboard-footer">
          <p>Incident history is stored locally and retained for 30 days</p>
        </footer>
      </div>
    </>
  );
}
