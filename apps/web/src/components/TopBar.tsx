import { Link, useLocation } from 'react-router-dom';
import type { WidgetResult } from '../types';

interface TopBarProps {
  searchQuery: string;
  onSearchChange: (query: string) => void;
  systemWidgets?: WidgetResult[];
}

export function TopBar({ searchQuery, onSearchChange, systemWidgets }: TopBarProps) {
  const location = useLocation();

  const getIcon = (label: string) => {
    switch(label) {
      case 'CPU': return <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><rect x="4" y="4" width="16" height="16" rx="0"/><path d="M9 9h6v6H9z"/><path d="M15 2v2M9 2v2M20 15h2M20 9h2M15 20v2M9 20v2M2 15h2M2 9h2"/></svg>;
      case 'RAM': return <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M6 19v2M10 19v2M14 19v2M18 19v2M8 11V9M12 11V9M16 11V9M20 11V9"/><rect x="2" y="5" width="20" height="14" rx="0"/></svg>;
      case 'Disk': return <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M22 12H2M5.45 5.11L2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11z"/><path d="M6 16h.01M10 16h.01"/></svg>;
      case 'Temp': return <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M14 4v10.54a4 4 0 1 1-4 0V4a2 2 0 0 1 4 0Z"/></svg>;
      case 'Uptime': return <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>;
      default: return null;
    }
  };

  return (
    <div className="topbar">
      <div className="topbar-inner">
        <div className="top-left">
          <Link to="/" className="brand-link">
            <img src="/logo.png" alt="homectl" className="logo-img" />
          </Link>
          <div className="nav-icons">
            <Link to="/" className={location.pathname === '/' ? 'active' : ''} title="Dashboard">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/>
              </svg>
            </Link>
            <Link to="/settings" className={location.pathname === '/settings' ? 'active' : ''} title="Settings">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"/><circle cx="12" cy="12" r="3"/>
              </svg>
            </Link>
          </div>
        </div>

        <div className="top-right">
          {systemWidgets && systemWidgets.length > 0 && (
            <div className="system-monitor-header">
              {systemWidgets.map((w, i) => {
                const isPercent = w.label === 'CPU' || w.label === 'RAM' || w.label === 'Disk';
                const val = typeof w.value === 'number' ? w.value : 0;
                return (
                  <div key={i} className={`sys-item ${w.state}`} title={`${w.label}: ${w.formatted}`}>
                    <span className="sys-icon">{getIcon(w.label || '')}</span>
                    {isPercent ? (
                      <div className="sys-bar-container">
                        <div className="sys-bar-fill" style={{ width: `${Math.min(100, val)}%` }} />
                      </div>
                    ) : (
                      <span className="sys-value">{w.formatted}</span>
                    )}
                  </div>
                );
              })}
              <div className="v-sep"></div>
            </div>
          )}
          <div className="search">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <circle cx="11" cy="11" r="8"/>
              <path d="m21 21-4.35-4.35"/>
            </svg>
            <input
              type="search"
              placeholder="Search services..."
              value={searchQuery}
              onChange={(e) => onSearchChange(e.target.value)}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
