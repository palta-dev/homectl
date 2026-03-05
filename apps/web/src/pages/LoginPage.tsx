import { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useLogin, useConfig } from '../hooks/api';
import { useAuthStore } from '../stores/authStore';

export function LoginPage() {
  const [password, setPassword] = useState('');
  const { mutate: login, isPending, error, isError } = useLogin();
  const { data: config } = useConfig();
  const { isAuthenticated } = useAuthStore();
  const navigate = useNavigate();
  const location = useLocation();

  const from = (location.state as { from?: { pathname: string } })?.from?.pathname || '/settings';

  useEffect(() => {
    if (isAuthenticated) {
      navigate(from, { replace: true });
    }
  }, [isAuthenticated, navigate, from]);

  // If no password is set on server, just redirect back
  useEffect(() => {
    if (config && !config.passwordProtected) {
      navigate('/', { replace: true });
    }
  }, [config, navigate]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    login(password);
  };

  return (
    <div className="main" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', minHeight: '60vh' }}>
      <div className="card" style={{ width: '100%', maxWidth: '400px', padding: '40px' }}>
        <div style={{ textAlign: 'center', marginBottom: '32px' }}>
          <img src="/logo.png" alt="homectl" style={{ height: '48px', marginBottom: '24px' }} />
          <h2 style={{ fontSize: '18px', fontWeight: 600, textTransform: 'uppercase', letterSpacing: '0.1em' }}>
            Protected Settings
          </h2>
          <p style={{ color: 'var(--muted)', fontSize: '14px', marginTop: '8px' }}>
            Please enter your password to continue
          </p>
        </div>

        <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
          <div className="field">
            <label>Password</label>
            <input 
              type="password" 
              value={password} 
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              autoFocus
              required
            />
          </div>

          {isError && (
            <div style={{ color: '#ef4444', fontSize: '13px', fontWeight: 500, textAlign: 'center' }}>
              {error.message || 'Invalid password'}
            </div>
          )}

          <button className="btn" type="submit" disabled={isPending}>
            {isPending ? 'Authenticating...' : 'Unlock Settings'}
          </button>
        </form>
      </div>
    </div>
  );
}
