import { useState, useEffect } from 'react';
import { TopBar } from '../components/TopBar';
import { useConfig, useUpdateConfig, useLogout } from '../hooks/api';
import { useSearchStore } from '../stores/searchStore';
import type { Group, Service } from '../types';

export function SettingsPage() {
  const { data: config } = useConfig();
  const updateConfig = useUpdateConfig();
  const logout = useLogout();
  const [searchQuery, setSearchQuery] = useState('');
  const { setQuery } = useSearchStore();
  const [activeTab, setActiveTab] = useState('general');

  const [isInitialized, setIsInitialized] = useState(false);
  const [title, setTitle] = useState('');
  const [theme, setTheme] = useState<'dark' | 'light'>('dark');
  const [background, setBackground] = useState('');
  const [backgroundOpacity, setBackgroundOpacity] = useState(0.5);
  const [allowHosts, setAllowHosts] = useState('');
  const [dockerEnabled, setDockerEnabled] = useState(false);
  const [dockerSubnet, setDockerSubnet] = useState('');
  const [dockerIgnore, setDockerIgnore] = useState('');
  const [requestTimeout, setRequestTimeout] = useState('10s');
  const [password, setPassword] = useState('');
  const [groups, setGroups] = useState<Group[]>([]);

  useEffect(() => {
    if (config && !isInitialized) {
      setTitle(config.title || '');
      setTheme(config.theme || 'dark');
      setBackground(config.background || '');
      setBackgroundOpacity(config.backgroundOpacity ?? 0.5);
      setAllowHosts(config.allowHosts?.join('\n') || '');
      setDockerEnabled(config.docker?.enabled || false);
      setDockerSubnet(config.docker?.subnet || '');
      setDockerIgnore(config.docker?.ignore?.join('\n') || '');
      setRequestTimeout(config.requestTimeout || '10s');
      setPassword(''); // Don't load password from server (it's hidden)
      setGroups(config.groups || []);
      setIsInitialized(true);
    }
  }, [config, isInitialized]);

  const handleSave = () => {
    updateConfig.mutate({
      settings: {
        title,
        theme,
        background,
        backgroundOpacity,
        allowHosts: allowHosts.split('\n').filter(h => h.trim() !== ''),
        requestTimeout,
        password: password || undefined, // Only send if changed
        groups,
        docker: {
          ...config?.docker,
          enabled: dockerEnabled,
          subnet: dockerSubnet,
          ignore: dockerIgnore.split('\n').filter(h => h.trim() !== ''),
        }
      }
    });
  };

  const handleLogout = () => {
    logout.mutate();
  };

  const addGroup = () => {
    setGroups([...groups, { name: 'New Group', layout: 'grid', collapsed: false, services: [] }]);
  };

  const removeGroup = (index: number) => {
    setGroups(groups.filter((_, i) => i !== index));
  };

  const updateGroupName = (index: number, name: string) => {
    const newGroups = [...groups];
    newGroups[index].name = name;
    setGroups(newGroups);
  };

  const addService = (groupIndex: number) => {
    const newGroups = [...groups];
    newGroups[groupIndex].services.push({ name: 'New Service', url: 'http://' });
    setGroups(newGroups);
  };

  const removeService = (groupIndex: number, serviceIndex: number) => {
    const newGroups = [...groups];
    newGroups[groupIndex].services = newGroups[groupIndex].services.filter((_, i) => i !== serviceIndex);
    setGroups(newGroups);
  };

  const updateService = (groupIndex: number, serviceIndex: number, field: keyof Service, value: string) => {
    const newGroups = [...groups];
    const service = newGroups[groupIndex].services[serviceIndex];
    if (field === 'name' || field === 'url' || field === 'icon' || field === 'description') {
      service[field] = value;
    }
    setGroups(newGroups);
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

      <div className="main">
        <section className="section">
          <div className="section-head">
            <h2>Configuration</h2>
            <div className="divider"></div>
          </div>

          <div className="settings-tabs">
            <button 
              className={`tab ${activeTab === 'general' ? 'active' : ''}`}
              onClick={() => setActiveTab('general')}
            >
              General
            </button>
            <button 
              className={`tab ${activeTab === 'services' ? 'active' : ''}`}
              onClick={() => setActiveTab('services')}
            >
              Services
            </button>
            <button 
              className={`tab ${activeTab === 'discovery' ? 'active' : ''}`}
              onClick={() => setActiveTab('discovery')}
            >
              Auto-Discovery
            </button>
            <button 
              className={`tab ${activeTab === 'theme' ? 'active' : ''}`}
              onClick={() => setActiveTab('theme')}
            >
              Theme
            </button>
            <button 
              className={`tab ${activeTab === 'security' ? 'active' : ''}`}
              onClick={() => setActiveTab('security')}
            >
              Security
            </button>
            <button 
              className={`tab ${activeTab === 'config' ? 'active' : ''}`}
              onClick={() => setActiveTab('config')}
            >
              Config File
            </button>
          </div>

          {activeTab === 'general' && (
            <div className="form-grid">
              <div className="field">
                <label>Dashboard Title</label>
                <input 
                  type="text" 
                  value={title} 
                  onChange={(e) => setTitle(e.target.value)} 
                />
              </div>
              <div className="field">
                <label>Default Theme</label>
                <select 
                  value={theme} 
                  onChange={(e) => setTheme(e.target.value as 'dark' | 'light')}
                >
                  <option value="dark">Dark</option>
                  <option value="light">Light</option>
                </select>
              </div>
              <div className="field" style={{ gridColumn: '1 / -1' }}>
                <label>Background Image URL</label>
                <input 
                  type="text" 
                  placeholder="e.g. https://images.unsplash.com/photo-..." 
                  value={background} 
                  onChange={(e) => setBackground(e.target.value)} 
                />
                <p className="hint">Direct link to an image (Unsplash, etc.)</p>
              </div>
              <div className="field" style={{ gridColumn: '1 / -1' }}>
                <label>Background Tint Level: {(backgroundOpacity * 100).toFixed(0)}%</label>
                <input 
                  type="range" 
                  min="0" 
                  max="1" 
                  step="0.05" 
                  value={backgroundOpacity} 
                  onChange={(e) => setBackgroundOpacity(parseFloat(e.target.value))} 
                  className="range-input"
                />
                <p className="hint">Darken the background image to make text more readable.</p>
              </div>
              <div className="field" style={{ gridColumn: '1 / -1' }}>
                <label>Allowed Hosts (for SSRF protection)</label>
                <textarea 
                  value={allowHosts} 
                  onChange={(e) => setAllowHosts(e.target.value)} 
                />
                <p className="hint">One host/CIDR per line. Only these hosts can be accessed by health checks and widgets.</p>
              </div>

              <div className="row-actions">
                <button 
                  className="btn" 
                  onClick={handleSave}
                  disabled={updateConfig.isPending}
                >
                  {updateConfig.isPending ? 'Saving...' : 'Save Changes'}
                </button>
                <button className="btn secondary" onClick={() => {
                   setIsInitialized(false);
                }}>Reset to Defaults</button>
              </div>
              {updateConfig.isSuccess && <p className="success-msg">Settings saved!</p>}
              {updateConfig.isError && <p className="error-msg">Failed to save settings: {updateConfig.error.message}</p>}
            </div>
          )}

          {activeTab === 'security' && (
            <div className="form-grid">
              <div className="field" style={{ gridColumn: '1 / -1' }}>
                <label>Change Settings Password</label>
                <input 
                  type="password" 
                  value={password} 
                  onChange={(e) => setPassword(e.target.value)} 
                  placeholder="Enter new password to secure settings"
                />
                <p className="hint">Protect your configuration. If set, this password will be required for any changes. Leave empty to disable protection.</p>
              </div>

              <div className="row-actions">
                <button 
                  className="btn" 
                  onClick={handleSave}
                  disabled={updateConfig.isPending}
                >
                  {updateConfig.isPending ? 'Updating...' : 'Update Password'}
                </button>
              </div>

              <div className="field" style={{ gridColumn: '1 / -1', borderTop: '1px solid var(--stroke)', paddingTop: '32px', marginTop: '32px' }}>
                <label>Current Session</label>
                <p className="hint" style={{ marginBottom: '16px' }}>You are currently authenticated. To end your session, click the button below.</p>
                <button className="btn secondary" onClick={handleLogout} style={{ width: 'fit-content' }}>
                  Log Out
                </button>
              </div>

              {updateConfig.isSuccess && <p className="success-msg">Security settings updated!</p>}
              {updateConfig.isError && <p className="error-msg">Error: {updateConfig.error.message}</p>}
            </div>
          )}

          {activeTab === 'services' && (
            <div className="services-config">
              <div className="hint" style={{ marginBottom: '24px' }}>
                Manage your dashboard groups and manual services here.
              </div>
              {groups.map((group, groupIndex) => (
                <div key={groupIndex} className="field" style={{ marginBottom: '32px', border: '1px solid var(--stroke2)', padding: '24px' }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '12px', flex: 1 }}>
                      <label style={{ margin: 0 }}>Group Name:</label>
                      <input 
                        type="text" 
                        value={group.name} 
                        onChange={(e) => updateGroupName(groupIndex, e.target.value)}
                        style={{ fontWeight: 'bold', border: 'none', background: 'transparent', padding: '0', fontSize: '16px' }}
                      />
                    </div>
                    <button className="btn secondary" onClick={() => removeGroup(groupIndex)} style={{ border: 'none', color: '#ef4444' }}>Delete Group</button>
                  </div>
                  
                  <div className="services-list" style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                    {group.services.map((svc, svcIndex) => (
                      <div key={svcIndex} style={{ display: 'flex', gap: '12px', padding: '16px', background: 'rgba(255,255,255,0.02)', border: '1px solid var(--stroke)', alignItems: 'center' }}>
                        <div style={{ flex: 1, display: 'grid', gridTemplateColumns: '1fr 2fr', gap: '16px' }}>
                          <input 
                            placeholder="Service Name" 
                            type="text" 
                            value={svc.name} 
                            onChange={(e) => updateService(groupIndex, svcIndex, 'name', e.target.value)}
                          />
                          <input 
                            placeholder="URL (e.g. http://localhost:8080)" 
                            type="text" 
                            value={svc.url} 
                            onChange={(e) => updateService(groupIndex, svcIndex, 'url', e.target.value)}
                          />
                        </div>
                        <button className="btn secondary" onClick={() => removeService(groupIndex, svcIndex)} style={{ padding: '0 12px', height: '42px', fontSize: '20px' }}>&times;</button>
                      </div>
                    ))}
                    <button className="btn secondary" onClick={() => addService(groupIndex)} style={{ marginTop: '8px', width: 'fit-content' }}>+ Add Service</button>
                  </div>
                </div>
              ))}

              <div className="row-actions">
                <button className="btn secondary" onClick={addGroup}>+ Create New Group</button>
                <button className="btn" onClick={handleSave} disabled={updateConfig.isPending}>
                   {updateConfig.isPending ? 'Saving...' : 'Save All Changes'}
                </button>
              </div>
              {updateConfig.isSuccess && <p className="success-msg">Services updated successfully!</p>}
              {updateConfig.isError && <p className="error-msg">Error: {updateConfig.error.message}</p>}
            </div>
          )}

          {activeTab === 'discovery' && (
            <div className="form-grid">
              <div className="field">
                <label>Enable Auto-Discovery</label>
                <select 
                  value={dockerEnabled.toString()} 
                  onChange={(e) => setDockerEnabled(e.target.value === 'true')}
                >
                  <option value="true">Enabled</option>
                  <option value="false">Disabled</option>
                </select>
              </div>
              <div className="field">
                <label>Request Timeout (e.g. 5s, 10s)</label>
                <input 
                  type="text" 
                  value={requestTimeout}
                  onChange={(e) => setRequestTimeout(e.target.value)}
                />
              </div>
              <div className="field" style={{ gridColumn: '1 / -1' }}>
                <label>Subnet Scan (optional)</label>
                <input 
                  type="text" 
                  placeholder="e.g., 192.168.1" 
                  value={dockerSubnet}
                  onChange={(e) => setDockerSubnet(e.target.value)}
                />
                <p className="hint">Enter subnet prefix to scan entire range (e.g., 192.168.1 scans 192.168.1.1-254)</p>
              </div>
              <div className="field" style={{ gridColumn: '1 / -1' }}>
                <label>Ignored Services</label>
                <textarea 
                  value={dockerIgnore} 
                  onChange={(e) => setDockerIgnore(e.target.value)} 
                  placeholder="e.g. MyService"
                />
                <p className="hint">One service name per line. These services will be hidden from auto-discovery results.</p>
              </div>

              <div className="row-actions">
                <button className="btn" onClick={handleSave} disabled={updateConfig.isPending}>
                   {updateConfig.isPending ? 'Saving...' : 'Save & Rescan'}
                </button>
              </div>
              {updateConfig.isSuccess && <p className="success-msg">Discovery settings updated!</p>}
              {updateConfig.isError && <p className="error-msg">Error: {updateConfig.error.message}</p>}
            </div>
          )}

          {activeTab === 'theme' && (
            <div className="form-grid">
              <div className="field">
                <label>Theme</label>
                <select 
                  value={theme} 
                  onChange={(e) => setTheme(e.target.value as 'dark' | 'light')}
                >
                  <option value="dark">Dark</option>
                  <option value="light">Light</option>
                </select>
              </div>
              <div className="field">
                <label>Background Image URL</label>
                <input 
                  type="text" 
                  placeholder="e.g. https://images.unsplash.com/photo-..." 
                  value={background} 
                  onChange={(e) => setBackground(e.target.value)} 
                />
              </div>
              <div className="field">
                <label>Background Tint Level: {(backgroundOpacity * 100).toFixed(0)}%</label>
                <input 
                  type="range" 
                  min="0" 
                  max="1" 
                  step="0.05" 
                  value={backgroundOpacity} 
                  onChange={(e) => setBackgroundOpacity(parseFloat(e.target.value))} 
                  className="range-input"
                />
              </div>

              <div className="row-actions">
                <button className="btn" onClick={handleSave} disabled={updateConfig.isPending}>
                  {updateConfig.isPending ? 'Applying...' : 'Apply Theme'}
                </button>
              </div>
              {updateConfig.isSuccess && <p className="success-msg">Theme updated!</p>}
            </div>
          )}

          {activeTab === 'config' && (
            <div className="field">
              <label>Current Configuration (config.yaml)</label>
              <pre className="code-block">
{`version: 1
settings:
  title: "${config?.title || 'homectl'}"
  theme: "${config?.theme || 'dark'}"
  allowHosts:
${config?.allowHosts?.map(h => `    - "${h}"`).join('\n') || ''}
  docker:
    enabled: ${config?.docker?.enabled || false}
    hosts:
${config?.docker?.hosts?.map(h => `      - address: "${h.address}"\n        tags: ${JSON.stringify(h.tags)}`).join('\n') || ''}
    ignore:
${config?.docker?.ignore?.map(i => `      - "${i}"`).join('\n') || ''}

groups:
${config?.groups?.map(g => `  - name: "${g.name}"\n    services:\n${g.services?.map(s => `      - name: "${s.name}"\n        url: "${s.url}"`).join('\n')}`).join('\n')}`}
              </pre>
              <div className="row-actions">
                <button className="btn" onClick={() => navigator.clipboard.writeText('config.yaml')}>
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                  </svg>
                  Copy Config
                </button>
                <a href="/config.yaml" download className="btn secondary">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                    <polyline points="7 10 12 15 17 10"/>
                    <line x1="12" y1="15" x2="12" y2="3"/>
                  </svg>
                  Download
                </a>
              </div>
            </div>
          )}
        </section>

        <footer className="dashboard-footer">
          <p>Changes to configuration require a server restart to take effect</p>
        </footer>
      </div>
    </>
  );
}
