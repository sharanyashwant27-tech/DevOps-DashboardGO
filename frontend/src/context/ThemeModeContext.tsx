import { createContext, useContext, useEffect, useMemo, useState, type ReactNode } from 'react';
import { createTheme, type Theme } from '@mui/material';

interface ThemeModeContextValue {
  mode: 'dark' | 'light';
  toggle: () => void;
  theme: Theme;
}

const ThemeModeContext = createContext<ThemeModeContextValue | undefined>(undefined);

export function ThemeModeProvider({ children }: { children: ReactNode }) {
  const [mode, setMode] = useState<'dark' | 'light'>(
    () => (localStorage.getItem('theme') as 'dark' | 'light') || 'dark',
  );

  useEffect(() => {
    localStorage.setItem('theme', mode);
    document.documentElement.classList.toggle('light', mode === 'light');
    document.documentElement.classList.toggle('dark', mode === 'dark');
  }, [mode]);

  const theme = useMemo(
    () =>
      createTheme({
        palette: {
          mode,
          primary: { main: '#0ea5a4' },
          secondary: { main: '#3b82f6' },
          background: {
            default: mode === 'dark' ? '#0b1220' : '#f4f7fb',
            paper: mode === 'dark' ? '#111827' : '#ffffff',
          },
        },
        typography: {
          fontFamily: '"IBM Plex Sans", system-ui, sans-serif',
          h1: { fontFamily: '"Sora", system-ui, sans-serif' },
          h2: { fontFamily: '"Sora", system-ui, sans-serif' },
          h3: { fontFamily: '"Sora", system-ui, sans-serif' },
          h4: { fontFamily: '"Sora", system-ui, sans-serif' },
          h5: { fontFamily: '"Sora", system-ui, sans-serif' },
          h6: { fontFamily: '"Sora", system-ui, sans-serif' },
        },
        shape: { borderRadius: 14 },
      }),
    [mode],
  );

  const value = useMemo(
    () => ({
      mode,
      toggle: () => setMode((m) => (m === 'dark' ? 'light' : 'dark')),
      theme,
    }),
    [mode, theme],
  );

  return <ThemeModeContext.Provider value={value}>{children}</ThemeModeContext.Provider>;
}

export function useThemeMode() {
  const ctx = useContext(ThemeModeContext);
  if (!ctx) throw new Error('useThemeMode must be used within ThemeModeProvider');
  return ctx;
}
