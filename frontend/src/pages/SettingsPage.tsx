import { Button, Chip, Paper, Typography } from '@mui/material';
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
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useThemeMode } from '../context/ThemeModeContext';
import StatCard from '../components/dashboard/StatCard';

const modules = [
  { title: 'Dashboard', subtitle: 'Live KPIs & charts', to: '/', icon: <DashboardIcon />, accent: '#0ea5a4' },
  { title: 'Projects', subtitle: 'Apps & organizations', to: '/projects', icon: <FolderIcon />, accent: '#3b82f6' },
  { title: 'Jenkins', subtitle: 'Pipelines & builds', to: '/jenkins', icon: <BuildIcon />, accent: '#d97706' },
  { title: 'GitHub', subtitle: 'Repos & Actions', to: '/github', icon: <GitHubIcon />, accent: '#8b5cf6' },
  { title: 'Docker', subtitle: 'Containers & images', to: '/docker', icon: <ViewInArIcon />, accent: '#06b6d4' },
  { title: 'Kubernetes', subtitle: 'Pods & deployments', to: '/kubernetes', icon: <HubIcon />, accent: '#6366f1' },
  { title: 'Servers', subtitle: 'Host metrics', to: '/servers', icon: <DnsIcon />, accent: '#14b8a6' },
  { title: 'Deployments', subtitle: 'Release history', to: '/deployments', icon: <RocketLaunchIcon />, accent: '#0ea5a4' },
  { title: 'Incidents', subtitle: 'Incident management', to: '/incidents', icon: <ReportIcon />, accent: '#f43f5e' },
  { title: 'Alerts', subtitle: 'Critical signals', to: '/alerts?severity=critical', icon: <NotificationsActiveIcon />, accent: '#f59e0b' },
];

export default function SettingsPage() {
  const { user } = useAuth();
  const { mode, toggle } = useThemeMode();

  function exportCSV() {
    const rows = [
      ['module', 'path'],
      ...modules.map((m) => [m.title, m.to]),
    ];
    const csv = rows.map((r) => r.join(',')).join('\n');
    const blob = new Blob([csv], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'devops-command-center-modules.csv';
    a.click();
  }

  return (
    <div className="space-y-6">
      <div>
        <Typography variant="h4" className="font-display">
          Settings
        </Typography>
        <Typography color="text.secondary">
          Operations Console preferences — click a module card to open its page
        </Typography>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-5 gap-4">
        {modules.map((m, i) => (
          <StatCard
            key={m.to}
            title={m.title}
            value="Open"
            subtitle={m.subtitle}
            icon={m.icon}
            accent={m.accent}
            to={m.to}
            delay={i * 30}
          />
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 max-w-5xl">
        <Paper className="glass-panel p-5 space-y-3" elevation={0}>
          <Typography variant="h6">Profile</Typography>
          <Typography>{user?.name}</Typography>
          <Typography color="text.secondary">{user?.email}</Typography>
          <Chip label={user?.role} color="primary" />
        </Paper>
        <Paper className="glass-panel p-5 space-y-3" elevation={0}>
          <Typography variant="h6">Appearance</Typography>
          <Typography color="text.secondary">Current mode: {mode}</Typography>
          <Button variant="outlined" onClick={toggle}>
            Toggle dark / light
          </Button>
        </Paper>
        <Paper className="glass-panel p-5 space-y-3" elevation={0}>
          <Typography variant="h6">Integrations</Typography>
          <Typography variant="body2" color="text.secondary">
            Configure connectors via environment variables / config.yaml, then open the matching module.
          </Typography>
          <div className="flex flex-wrap gap-2">
            <Button component={Link} to="/jenkins" size="small" variant="outlined">
              Jenkins
            </Button>
            <Button component={Link} to="/github" size="small" variant="outlined">
              GitHub
            </Button>
            <Button component={Link} to="/docker" size="small" variant="outlined">
              Docker
            </Button>
            <Button component={Link} to="/kubernetes" size="small" variant="outlined">
              Kubernetes
            </Button>
            <Button component={Link} to="/alerts" size="small" variant="outlined">
              Alerts
            </Button>
          </div>
          <ul className="list-disc pl-5 text-sm space-y-1 opacity-80">
            <li>DCC_JENKINS_URL / DCC_JENKINS_TOKEN</li>
            <li>DCC_GITHUB_TOKEN</li>
            <li>DCC_SLACK_WEBHOOK_URL</li>
            <li>DCC_TEAMS_WEBHOOK_URL</li>
          </ul>
        </Paper>
        <Paper className="glass-panel p-5 space-y-3" elevation={0}>
          <Typography variant="h6">Exports</Typography>
          <Button variant="contained" onClick={exportCSV}>
            Export modules CSV
          </Button>
        </Paper>
      </div>
    </div>
  );
}
