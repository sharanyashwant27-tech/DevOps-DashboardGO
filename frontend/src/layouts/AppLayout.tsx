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

const width = 260;

const nav = [
  { to: '/', label: 'Dashboard', icon: <DashboardIcon /> },
  { to: '/projects', label: 'Projects', icon: <FolderIcon /> },
  { to: '/jenkins', label: 'Jenkins', icon: <BuildIcon /> },
  { to: '/github', label: 'GitHub', icon: <GitHubIcon /> },
  { to: '/docker', label: 'Docker', icon: <ViewInArIcon /> },
  { to: '/kubernetes', label: 'Kubernetes', icon: <HubIcon /> },
  { to: '/servers', label: 'Servers', icon: <DnsIcon /> },
  { to: '/deployments', label: 'Deployments', icon: <RocketLaunchIcon /> },
  { to: '/incidents', label: 'Incidents', icon: <ReportIcon /> },
  { to: '/alerts', label: 'Alerts', icon: <NotificationsActiveIcon /> },
  { to: '/settings', label: 'Settings', icon: <SettingsIcon /> },
];

export default function AppLayout() {
  const { user, logout } = useAuth();
  const { mode, toggle } = useThemeMode();
  const navigate = useNavigate();
  const isMobile = useMediaQuery('(max-width:900px)');
  const [open, setOpen] = useState(!isMobile);

  const drawer = (
    <Box className="h-full flex flex-col" sx={{ bgcolor: 'background.paper' }}>
      <Box className="px-5 py-6">
        <Typography variant="h6" className="font-display tracking-tight">
          DevOps Command Center
        </Typography>
        <Typography variant="caption" color="text.secondary">
          Multi-cloud operations hub
        </Typography>
      </Box>
      <List className="flex-1 px-2">
        {nav.map((item) => (
          <ListItemButton
            key={item.to}
            component={NavLink}
            to={item.to}
            end={item.to === '/'}
            sx={{
              borderRadius: 2,
              mb: 0.5,
              '&.active': { bgcolor: 'rgba(14,165,164,0.15)', color: 'primary.main' },
            }}
          >
            <ListItemIcon sx={{ minWidth: 40, color: 'inherit' }}>{item.icon}</ListItemIcon>
            <ListItemText primary={item.label} />
          </ListItemButton>
        ))}
      </List>
      <Box className="p-4 border-t border-white/10">
        <Typography variant="body2" className="font-medium">
          {user?.name}
        </Typography>
        <Chip size="small" label={user?.role} sx={{ mt: 1 }} color="primary" variant="outlined" />
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
          [`& .MuiDrawer-paper`]: { width, boxSizing: 'border-box', borderRight: '1px solid var(--line)' },
        }}
      >
        {drawer}
      </Drawer>
      <Box className="flex-1 min-w-0">
        <AppBar
          position="sticky"
          elevation={0}
          color="transparent"
          sx={{ borderBottom: '1px solid var(--line)', backdropFilter: 'blur(10px)' }}
        >
          <Toolbar className="gap-2">
            {isMobile && (
              <IconButton onClick={() => setOpen(true)}>
                <MenuIcon />
              </IconButton>
            )}
            <Typography variant="h6" className="flex-1 font-display">
              Operations Console
            </Typography>
            <IconButton onClick={toggle} aria-label="toggle theme">
              {mode === 'dark' ? <LightModeIcon /> : <DarkModeIcon />}
            </IconButton>
            <IconButton
              onClick={() => {
                logout();
                navigate('/login');
              }}
              aria-label="logout"
            >
              <LogoutIcon />
            </IconButton>
            <Link to="/settings" className="text-sm opacity-80 hover:opacity-100">
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
