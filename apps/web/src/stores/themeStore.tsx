import { create } from 'zustand';

type Theme = 'dark';

interface ThemeState {
  theme: Theme;
  toggleTheme: () => void;
  setTheme: (theme: Theme) => void;
}

export const useThemeStore = create<ThemeState>(() => ({
  theme: 'dark',
  toggleTheme: () => {
    // Dark mode only - no-op
    console.log('Dark mode only');
  },
  setTheme: () => {
    // Dark mode only - no-op
  },
}));

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  return <>{children}</>;
}
