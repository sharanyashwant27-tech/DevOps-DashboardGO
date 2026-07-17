import { useState } from 'react';
import { Link, NavLink, Outlet, useNavigate } from 'react-router-dom';
import {
  AppBar,
  Box,
  Drawer,
  IconButton,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Toolbar,
  Typography,
  Chip,
  useMediaQuery,
} from '@mui/material';
import MenuIcon from '@mui/icons-material/Menu';
import DashboardIcon from '@mui/icons-material/Dashboard';
import FolderIcon from '@mui/icons-material/Folder';
import BuildIcon from '@mui/icons-material/Build';
import GitHubIcon from '@mui/icons-material/GitHub';
import ViewInArIcon from '@mui/icons-material/ViewInAr';
import HubIcon from '@mui/icons-material/Hub';
import DnsIcon from '@mui/icons-material/Dns';
import RocketLaunchIcon from '@mui/icons-material/RocketLaunch';
import ReportIcon from '@mui/icons-material/Report';
import NotificationsActiveIcon from '@mui/icons-material/NotificationsActive';
import SettingsIcon from '@mui/icons-material/Settings';
import DarkModeIcon from '@mui/icons-material/DarkMode';
import LightModeIcon from '@mui/icons-material/LightMode';
import LogoutIcon from '@mui/icons-material/Logout';
import { useAuth } from '../context/AuthContext';
import { useThemeMode } from '../context/ThemeModeContext';

const width = 268;

const nav = [
  { to: '/', label: 'Dashboard', icon: <DashboardIcon />, accent: '#22d3ee' },
  { to: '/projects', label: 'Projects', icon: <FolderIcon />, accent: '#38bdf8' },
  { to: '/jenkins', label: 'Jenkins', icon: <BuildIcon />, accent: '#fbbf24' },
  { to: '/github', label: 'GitHub', icon: <GitHubIcon />, accent: '#a78bfa' },
  { to: '/docker', label: 'Docker', icon: <ViewInArIcon />, accent: '#06b6d4' },
  { to: '/kubernetes', label: 'Kubernetes', icon: <HubIcon />, accent: '#60a5fa' },
  { to: '/servers', label: 'Servers', icon: <DnsIcon />, accent: '#34d399' },
  { to: '/deployments', label: 'Deployments', icon: <RocketLaunchIcon />, accent: '#2dd4bf' },
  { to: '/incidents', label: 'Incidents', icon: <ReportIcon />, accent: '#fb7185' },
  { to: '/alerts', label: 'Alerts', icon: <NotificationsActiveIcon />, accent: '#f59e0b' },
  { to: '/settings', label: 'Settings', icon: <SettingsIcon />, accent: '#94a3b8' },
];

export default function AppLayout() {
  const { user, logout } = useAuth();
  const { mode, toggle } = useThemeMode();
  const navigate = useNavigate();
  const isMobile = useMediaQuery('(max-width:900px)');
  const [open, setOpen] = useState(!isMobile);

  const drawer = (
    <Box className="h-full flex flex-col nav-rail">
      <Box className="px-5 py-6">
        <div className="flex items-center gap-3 mb-2">
          <span className="brand-mark">DC</span>
          <div>
            <Typography variant="subtitle1" className="font-display leading-tight">
              DevOps Command
            </Typography>
            <Typography variant="caption" sx={{ color: 'text.secondary' }}>
              Colorful ops console
            </Typography>
          </div>
        </div>
      </Box>
      <List className="flex-1 px-2 overflow-auto">
        {nav.map((item) => (
          <ListItemButton
            key={item.to}
            component={NavLink}
            to={item.to}
            end={item.to === '/'}
            sx={{
              borderRadius: 2.5,
              mb: 0.5,
              transition: 'all 0.2s ease',
              '&:hover': {
                bgcolor: `${item.accent}18`,
              },
              '&.active': {
                bgcolor: `${item.accent}22`,
                color: item.accent,
                boxShadow: `inset 3px 0 0 ${item.accent}`,
              },
            }}
          >
            <ListItemIcon sx={{ minWidth: 40, color: 'inherit' }}>{item.icon}</ListItemIcon>
            <ListItemText primary={item.label} primaryTypographyProps={{ fontWeight: 600 }} />
          </ListItemButton>
        ))}
      </List>
      <Box
        className="p-4 m-3 rounded-2xl"
        sx={{
          border: '1px solid var(--line)',
          background: 'linear-gradient(135deg, rgba(34,211,238,0.12), rgba(251,191,36,0.08))',
        }}
      >
        <Typography variant="body2" className="font-medium">
          {user?.name}
        </Typography>
        <Chip
          size="small"
          label={user?.role}
          sx={{
            mt: 1,
            bgcolor: 'rgba(34,211,238,0.18)',
            color: 'primary.main',
            border: '1px solid rgba(34,211,238,0.35)',
          }}
        />
      </Box>
    </Box>
  );

  return (
    <Box className="min-h-screen flex">
      <Drawer
        variant={isMobile ? 'temporary' : 'permanent'}
        open={open}
        onClose={() => setOpen(false)}
        sx={{
          width,
          [`& .MuiDrawer-paper`]: {
            width,
            boxSizing: 'border-box',
            borderRight: '1px solid var(--line)',
            background: 'transparent',
          },
        }}
      >
        {drawer}
      </Drawer>
      <Box className="flex-1 min-w-0">
        <AppBar
          position="sticky"
          elevation={0}
          color="transparent"
          sx={{
            borderBottom: '1px solid var(--line)',
            backdropFilter: 'blur(14px)',
            background: 'linear-gradient(90deg, rgba(34,211,238,0.08), rgba(251,191,36,0.06), transparent 70%)',
          }}
        >
          <Toolbar className="gap-2">
            {isMobile && (
              <IconButton onClick={() => setOpen(true)} sx={{ color: 'primary.main' }}>
                <MenuIcon />
              </IconButton>
            )}
            <Typography variant="h6" className="flex-1 font-display">
              Operations Console
            </Typography>
            <Chip
              size="small"
              label="Live stack"
              sx={{
                display: { xs: 'none', sm: 'inline-flex' },
                bgcolor: 'rgba(52,211,153,0.15)',
                color: '#34d399',
                border: '1px solid rgba(52,211,153,0.35)',
              }}
            />
            <IconButton onClick={toggle} aria-label="toggle theme" sx={{ color: 'secondary.main' }}>
              {mode === 'dark' ? <LightModeIcon /> : <DarkModeIcon />}
            </IconButton>
            <IconButton
              onClick={() => {
                logout();
                navigate('/login');
              }}
              aria-label="logout"
              sx={{ color: 'error.main' }}
            >
              <LogoutIcon />
            </IconButton>
            <Link to="/settings" className="link-accent text-sm">
              Settings
            </Link>
          </Toolbar>
        </AppBar>
        <Box className="p-4 md:p-8 page-enter">
          <Outlet />
        </Box>
      </Box>
    </Box>
  );
}
