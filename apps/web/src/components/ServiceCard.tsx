import { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { cn, formatLatency } from '../lib/utils';
import type { ServiceWithStatus } from '../types';
import { useConfig, useUpdateConfig } from '../hooks/api';
import { useQueryClient } from '@tanstack/react-query';

interface ServiceCardProps {
  service: ServiceWithStatus;
}

export function ServiceCard({ service }: ServiceCardProps) {
  const { name, url, favicon, description, status, tags, newTab, isDiscovered } = service;
  if (service.widgets && service.widgets.length > 0) {
    console.log(`Service ${name} has widgets:`, service.widgets);
  }
  const [showMenu, setShowMenu] = useState(false);
  const [showDetails, setShowDetails] = useState(false);
  
  const menuRef = useRef<HTMLDivElement>(null);
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { data: config } = useConfig();
  const updateConfig = useUpdateConfig();

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setShowMenu(false);
      }
    }
    if (showMenu) {
      document.addEventListener('mousedown', handleClickOutside);
    }
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [showMenu]);

  const handleHide = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (!config) return;

    const currentIgnore = config.docker?.ignore || [];
    if (!currentIgnore.includes(name)) {
      updateConfig.mutate({
        settings: {
          docker: {
            ...config.docker,
            enabled: config.docker?.enabled ?? false,
            ignore: [...currentIgnore, name]
          }
        }
      });
    }
    setShowMenu(false);
  };

  const handleCopy = (e: React.MouseEvent) => {
    e.stopPropagation();
    navigator.clipboard.writeText(url);
    setShowMenu(false);
  };

  const handleCheckNow = (e: React.MouseEvent) => {
    e.stopPropagation();
    // Force a refresh of the services query
    queryClient.invalidateQueries({ queryKey: ['services'] });
    setTimeout(() => {
      setShowMenu(false);
    }, 1000);
  };

  const handleEdit = (e: React.MouseEvent) => {
    e.stopPropagation();
    navigate('/settings');
  };

  const toggleDetails = (e: React.MouseEvent) => {
    e.stopPropagation();
    setShowDetails(!showDetails);
    setShowMenu(false);
  };

  return (
    <article className={cn('card', showMenu && 'has-menu')} onClick={() => window.open(url, newTab ? '_blank' : '_self')}>
      <div className="shine"></div>
      <div className="card-inner">
        <div className="card-top">
          <div className="svc">
            <div className="svc-ico">
              {favicon ? (
                <img src={favicon} alt={name} className="favicon" onError={(e) => {
                  (e.target as HTMLImageElement).style.display = 'none';
                  (e.target as HTMLImageElement).nextElementSibling?.classList.add('show');
                }} />
              ) : null}
              <div className={cn('svc-fallback', favicon ? 'hidden' : '')}>
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
                  <polyline points="3.29 7 12 12 20.71 7"/>
                  <line x1="12" y1="22" x2="12" y2="12"/>
                </svg>
              </div>
            </div>
            <div className="svc-meta">
              <div className="svc-name">{name}</div>
              {description && <div className="svc-sub">{description}</div>}
              {tags && tags.length > 0 && (
                <div className="svc-sub">
                  {tags.slice(0, 3).map((tag, i) => (
                    <span key={i} className="pilltag">{tag}</span>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>

        {status.latency !== undefined && (
          <div className="metrics">
            <div>
              <div className="lat">{formatLatency(status.latency)}</div>
              <div className="sub">Response time</div>
            </div>
            <div className="spark">
              <svg viewBox="0 0 170 66" preserveAspectRatio="none">
                <defs>
                  <linearGradient id="grad" x1="0%" y1="0%" x2="0%" y2="100%">
                    <stop offset="0%" stopColor="var(--green)" stopOpacity="0.3"/>
                    <stop offset="100%" stopColor="var(--green)" stopOpacity="0"/>
                  </linearGradient>
                </defs>
                <path
                  d="M0,60 Q20,55 40,58 T80,52 T120,56 T170,50 L170,66 L0,66 Z"
                  fill="url(#grad)"
                />
                <path
                  d="M0,60 Q20,55 40,58 T80,52 T120,56 T170,50"
                  fill="none"
                  stroke="var(--green)"
                  strokeWidth="2"
                />
              </svg>
              <div className="label">24h</div>
            </div>
          </div>
        )}

        {showDetails && (
          <div className="details-overlay" onClick={e => e.stopPropagation()}>
            <div className="details-row">
              <span className="details-label">Endpoint</span>
              <span className="details-value">{url}</span>
            </div>
            <div className="details-row">
              <span className="details-label">Type</span>
              <span className="details-value">{isDiscovered ? 'Auto-Discovered' : 'Static Config'}</span>
            </div>
            {status.lastCheck && (
              <div className="details-row">
                <span className="details-label">Last Check</span>
                <span className="details-value">{new Date(status.lastCheck).toLocaleTimeString()}</span>
              </div>
            )}
          </div>
        )}

        {service.widgets && service.widgets.length > 0 && (
          <div className="card-widgets">
            {service.widgets.map((widget, i) => (
              <div key={i} className={cn('widget-pill', widget.state)}>
                <span className="widget-label">{widget.label}</span>
                <span className="widget-value">{widget.formatted || String(widget.value)}</span>
              </div>
            ))}
          </div>
        )}

        <div className="actions">
          <div className="act">
            <button onClick={handleCopy}>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
              </svg>
              Copy URL
            </button>
            <button onClick={(e) => { e.stopPropagation(); window.open(url, '_blank'); }}>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
                <polyline points="15 3 21 3 21 9"/>
                <line x1="10" y1="14" x2="21" y2="3"/>
              </svg>
              Open
            </button>
          </div>
          <div className="more-container" ref={menuRef}>
            <button className="more-btn" onClick={(e) => { e.stopPropagation(); setShowMenu(!showMenu); }}>
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="12" cy="12" r="1.5"/>
                <circle cx="12" cy="5" r="1.5"/>
                <circle cx="12" cy="19" r="1.5"/>
              </svg>
            </button>
            {showMenu && (
              <div className="dropdown-menu">
                <div className="menu-section">
                  <button onClick={handleCheckNow}>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M23 4v6h-6M1 20v-6h6M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/></svg>
                    Refresh Status
                  </button>
                  <button onClick={toggleDetails}>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
                    {showDetails ? 'Hide' : 'Show'} Details
                  </button>
                </div>
                
                <div className="menu-divider"></div>

                <div className="menu-section">
                  <button onClick={handleCopy}>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
                    Copy URL
                  </button>
                </div>

                <div className="menu-divider"></div>

                <div className="menu-section">
                  <button onClick={handleEdit}>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
                    Configure
                  </button>
                  {isDiscovered && (
                    <button onClick={handleHide} className="danger">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg>
                      Hide Service
                    </button>
                  )}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </article>
  );
}


