import { useQuery } from '@tanstack/react-query';
import { Chip, Typography } from '@mui/material';
import FolderIcon from '@mui/icons-material/Folder';
import PlayCircleIcon from '@mui/icons-material/PlayCircle';
import ErrorIcon from '@mui/icons-material/Error';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import DnsIcon from '@mui/icons-material/Dns';
import ViewInArIcon from '@mui/icons-material/ViewInAr';
import HubIcon from '@mui/icons-material/Hub';
import WarningAmberIcon from '@mui/icons-material/WarningAmber';
import RocketLaunchIcon from '@mui/icons-material/RocketLaunch';
import ReportIcon from '@mui/icons-material/Report';
import MemoryIcon from '@mui/icons-material/Memory';
import StorageIcon from '@mui/icons-material/Storage';
import SpeedIcon from '@mui/icons-material/Speed';
import StatCard from '../components/dashboard/StatCard';
import UsageCharts from '../components/metrics/UsageCharts';
import { dashboardApi } from '../services/api';
import { useWebSocket } from '../hooks/useWebSocket';
import { useCallback, useState } from 'react';
import { Link } from 'react-router-dom';
import type { DashboardStats } from '../types';

export default function DashboardPage() {
  const [live, setLive] = useState<DashboardStats | null>(null);
  const { data, isLoading } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: async () => (await dashboardApi.stats()).data.data!,
    refetchInterval: 15000,
  });

  const onMessage = useCallback((msg: unknown) => {
    const m = msg as { channel?: string; payload?: DashboardStats };
    if (m.channel === 'dashboard' && m.payload) setLive(m.payload);
  }, []);
  const { connected } = useWebSocket(onMessage);

  const stats = live || data;

  return (
    <div className="space-y-6">
      <div className="page-hero flex items-center justify-between gap-3 flex-wrap">
        <div className="relative z-[1]">
          <Typography variant="overline" sx={{ color: 'primary.main', letterSpacing: 2 }}>
            Operations overview
          </Typography>
          <Typography variant="h4" className="font-display">
            Command Dashboard
          </Typography>
          <Typography color="text.secondary" className="max-w-xl">
            Click any card to open its module — live color-coded view across CI/CD, clusters, and incidents
          </Typography>
        </div>
        <Chip
          label={connected ? 'Live WebSocket' : 'Polling'}
          className={connected ? 'live-chip animate-pulseSoft' : ''}
          color={connected ? 'success' : 'default'}
          variant="outlined"
        />
      </div>

      {isLoading && !stats ? (
        <Typography>Loading metrics...</Typography>
      ) : (
        <>
          <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-5 gap-4">
            <StatCard title="Projects" value={stats?.total_projects ?? 0} icon={<FolderIcon />} accent="#38bdf8" delay={0} to="/projects" />
            <StatCard
              title="Running Builds"
              value={stats?.running_builds ?? 0}
              icon={<PlayCircleIcon />}
              accent="#60a5fa"
              delay={40}
              to="/jenkins?filter=running"
            />
            <StatCard
              title="Failed Builds"
              value={stats?.failed_builds ?? 0}
              icon={<ErrorIcon />}
              accent="#fb7185"
              delay={80}
              to="/jenkins?filter=failed"
            />
            <StatCard
              title="Successful Builds"
              value={stats?.successful_builds ?? 0}
              icon={<CheckCircleIcon />}
              accent="#34d399"
              delay={120}
              to="/jenkins?filter=success"
            />
            <StatCard title="Servers Online" value={stats?.servers_online ?? 0} icon={<DnsIcon />} accent="#2dd4bf" delay={160} to="/servers" />
            <StatCard
              title="Containers"
              value={stats?.docker_containers_running ?? 0}
              icon={<ViewInArIcon />}
              accent="#22d3ee"
              delay={200}
              to="/docker"
            />
            <StatCard title="Pods Running" value={stats?.pods_running ?? 0} icon={<HubIcon />} accent="#818cf8" delay={240} to="/kubernetes" />
            <StatCard
              title="Critical Alerts"
              value={stats?.critical_alerts ?? 0}
              icon={<WarningAmberIcon />}
              accent="#fbbf24"
              delay={280}
              to="/alerts?severity=critical&status=open"
            />
            <StatCard
              title="Deployments Today"
              value={stats?.deployments_today ?? 0}
              icon={<RocketLaunchIcon />}
              accent="#a3e635"
              delay={320}
              to="/deployments"
            />
            <StatCard title="Incidents" value={stats?.open_incidents ?? 0} icon={<ReportIcon />} accent="#f43f5e" delay={360} to="/incidents" />
          </div>

          <UsageCharts
            cpu={stats?.cpu_usage ?? 0}
            memory={stats?.memory_usage ?? 0}
            disk={stats?.disk_usage ?? 0}
            successBuilds={stats?.successful_builds ?? 0}
            failedBuilds={stats?.failed_builds ?? 0}
          />

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <Link to="/servers" className="metric-tile metric-tile-cyan block no-underline text-inherit">
              <div className="flex items-center justify-between">
                <Typography variant="overline">CPU</Typography>
                <SpeedIcon fontSize="small" sx={{ color: '#22d3ee' }} />
              </div>
              <Typography variant="h3" className="font-display" sx={{ color: '#22d3ee' }}>
                {(stats?.cpu_usage ?? 0).toFixed(1)}%
              </Typography>
              <Typography variant="caption" sx={{ color: '#67e8f9' }}>
                Open servers →
              </Typography>
            </Link>
            <Link to="/servers" className="metric-tile metric-tile-sky block no-underline text-inherit">
              <div className="flex items-center justify-between">
                <Typography variant="overline">Memory</Typography>
                <MemoryIcon fontSize="small" sx={{ color: '#38bdf8' }} />
              </div>
              <Typography variant="h3" className="font-display" sx={{ color: '#38bdf8' }}>
                {(stats?.memory_usage ?? 0).toFixed(1)}%
              </Typography>
              <Typography variant="caption" sx={{ color: '#7dd3fc' }}>
                Open servers →
              </Typography>
            </Link>
            <Link to="/servers" className="metric-tile metric-tile-amber block no-underline text-inherit">
              <div className="flex items-center justify-between">
                <Typography variant="overline">Disk</Typography>
                <StorageIcon fontSize="small" sx={{ color: '#fbbf24' }} />
              </div>
              <Typography variant="h3" className="font-display" sx={{ color: '#fbbf24' }}>
                {(stats?.disk_usage ?? 0).toFixed(1)}%
              </Typography>
              <Typography variant="caption" sx={{ color: '#fcd34d' }}>
                Open servers →
              </Typography>
            </Link>
          </div>
        </>
      )}
    </div>
  );
}
