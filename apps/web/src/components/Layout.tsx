import { useMemo, useEffect } from 'react';
import { Outlet } from 'react-router-dom';
import { useConfig } from '../hooks/api';

export function Layout() {
  const { data: config } = useConfig();

  useEffect(() => {
    const theme = config?.theme || 'dark';
    document.documentElement.setAttribute('data-theme', theme);
  }, [config?.theme]);

  const style = useMemo(() => {
    if (!config?.background) return {};
    // Use 0.5 as default if backgroundOpacity is undefined or null
    const opacity = config.backgroundOpacity ?? 0.5;
    
    return {
      backgroundImage: `linear-gradient(rgba(0, 0, 0, ${opacity}), rgba(0, 0, 0, ${opacity})), url(${config.background})`,
      backgroundSize: 'cover',
      backgroundPosition: 'center',
      backgroundAttachment: 'fixed',
      backgroundRepeat: 'no-repeat',
      backgroundColor: 'var(--bg0)', // Fallback color
    };
  }, [config?.background, config?.backgroundOpacity]);

  return (
    <div className="app" style={style}>
      <main className="main">
        <Outlet />
      </main>
    </div>
  );
}
