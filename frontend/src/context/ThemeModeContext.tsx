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
          primary: {
            main: mode === 'dark' ? '#22d3ee' : '#0891b2',
            light: mode === 'dark' ? '#67e8f9' : '#22d3ee',
            dark: mode === 'dark' ? '#0891b2' : '#0e7490',
            contrastText: mode === 'dark' ? '#041018' : '#ffffff',
          },
          secondary: {
            main: mode === 'dark' ? '#fbbf24' : '#d97706',
            contrastText: '#1a1200',
          },
          success: { main: '#34d399' },
          warning: { main: '#fbbf24' },
          error: { main: '#fb7185' },
          info: { main: '#38bdf8' },
          background: {
            default: mode === 'dark' ? '#07101f' : '#e8f0fb',
            paper: mode === 'dark' ? '#12243a' : '#ffffff',
          },
          text: {
            primary: mode === 'dark' ? '#e8f1ff' : '#0b1f36',
            secondary: mode === 'dark' ? '#94a8c7' : '#5b6f8a',
          },
          divider: mode === 'dark' ? 'rgba(125, 211, 252, 0.16)' : 'rgba(15, 60, 110, 0.1)',
        },
        typography: {
          fontFamily: '"IBM Plex Sans", system-ui, sans-serif',
          h1: { fontFamily: '"Sora", system-ui, sans-serif', fontWeight: 700 },
          h2: { fontFamily: '"Sora", system-ui, sans-serif', fontWeight: 700 },
          h3: { fontFamily: '"Sora", system-ui, sans-serif', fontWeight: 700 },
          h4: { fontFamily: '"Sora", system-ui, sans-serif', fontWeight: 700 },
          h5: { fontFamily: '"Sora", system-ui, sans-serif', fontWeight: 600 },
          h6: { fontFamily: '"Sora", system-ui, sans-serif', fontWeight: 600 },
          button: { textTransform: 'none', fontWeight: 600 },
        },
        shape: { borderRadius: 14 },
        components: {
          MuiCssBaseline: {
            styleOverrides: {
              body: { backgroundColor: 'transparent' },
            },
          },
          MuiButton: {
            styleOverrides: {
              root: {
                borderRadius: 12,
                boxShadow: 'none',
              },
              containedPrimary: {
                backgroundImage: 'linear-gradient(135deg, #22d3ee, #38bdf8 60%, #67e8f9)',
                color: '#041018',
                '&:hover': {
                  backgroundImage: 'linear-gradient(135deg, #06b6d4, #0ea5e9 60%, #38bdf8)',
                  boxShadow: '0 8px 24px rgba(34, 211, 238, 0.35)',
                },
              },
              outlined: {
                borderWidth: 1.5,
              },
            },
          },
          MuiPaper: {
            styleOverrides: {
              root: {
                backgroundImage: 'none',
              },
            },
          },
          MuiChip: {
            styleOverrides: {
              root: {
                fontWeight: 600,
              },
              colorPrimary: {
                backgroundColor: mode === 'dark' ? 'rgba(34, 211, 238, 0.16)' : 'rgba(8, 145, 178, 0.12)',
              },
            },
          },
          MuiTextField: {
            defaultProps: {
              variant: 'outlined',
            },
            styleOverrides: {
              root: {
                '& .MuiOutlinedInput-root': {
                  backgroundColor: mode === 'dark' ? 'rgba(7, 16, 31, 0.45)' : 'rgba(255,255,255,0.7)',
                  borderRadius: 12,
                },
              },
            },
          },
          MuiAppBar: {
            styleOverrides: {
              root: {
                backgroundImage: 'none',
              },
            },
          },
          MuiDrawer: {
            styleOverrides: {
              paper: {
                backgroundImage: 'none',
                borderRight: '1px solid var(--line)',
              },
            },
          },
          MuiTableCell: {
            styleOverrides: {
              root: {
                borderColor: 'var(--line)',
              },
            },
          },
        },
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
