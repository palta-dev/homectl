import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface AuthState {
  password: string | null;
  isAuthenticated: boolean;
  setPassword: (password: string | null) => void;
  setAuthenticated: (auth: boolean) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      password: null,
      isAuthenticated: false,
      setPassword: (password) => set({ password }),
      setAuthenticated: (isAuthenticated) => set({ isAuthenticated }),
      logout: () => set({ password: null, isAuthenticated: false }),
    }),
    {
      name: 'homectl-auth',
    }
  )
);
