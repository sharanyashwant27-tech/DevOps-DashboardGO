import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from 'chart.js';
import { Line, Bar } from 'react-chartjs-2';
import { Paper, Typography } from '@mui/material';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend,
  Filler,
);

interface Props {
  cpu: number;
  memory: number;
  disk: number;
  successBuilds: number;
  failedBuilds: number;
}

export default function UsageCharts({ cpu, memory, disk, successBuilds, failedBuilds }: Props) {
  const labels = ['-50m', '-40m', '-30m', '-20m', '-10m', 'now'];
  const trend = (base: number) =>
    labels.map((_, i) => Math.max(0, Math.min(100, base + Math.sin(i) * 8 - 4 + i)));

  const lineData = {
    labels,
    datasets: [
      {
        label: 'CPU %',
        data: trend(cpu),
        borderColor: '#0ea5a4',
        backgroundColor: 'rgba(14,165,164,0.15)',
        fill: true,
        tension: 0.35,
      },
      {
        label: 'Memory %',
        data: trend(memory),
        borderColor: '#3b82f6',
        backgroundColor: 'rgba(59,130,246,0.12)',
        fill: true,
        tension: 0.35,
      },
      {
        label: 'Disk %',
        data: trend(disk),
        borderColor: '#f59e0b',
        backgroundColor: 'rgba(245,158,11,0.1)',
        fill: true,
        tension: 0.35,
      },
    ],
  };

  const barData = {
    labels: ['Successful', 'Failed'],
    datasets: [
      {
        label: 'Builds',
        data: [successBuilds, failedBuilds],
        backgroundColor: ['#0ea5a4', '#ef4444'],
        borderRadius: 8,
      },
    ],
  };

  return (
    <div className="grid grid-cols-1 xl:grid-cols-3 gap-4">
      <Paper className="glass-panel p-4 xl:col-span-2" elevation={0}>
        <Typography variant="h6" className="mb-3 font-display">
          Resource Usage
        </Typography>
        <Line
          data={lineData}
          options={{
            responsive: true,
            plugins: { legend: { labels: { color: '#94a3b8' } } },
            scales: {
              x: { ticks: { color: '#94a3b8' }, grid: { color: 'rgba(148,163,184,0.1)' } },
              y: { ticks: { color: '#94a3b8' }, grid: { color: 'rgba(148,163,184,0.1)' }, max: 100 },
            },
          }}
        />
      </Paper>
      <Paper className="glass-panel p-4" elevation={0}>
        <Typography variant="h6" className="mb-3 font-display">
          Build Success Rate
        </Typography>
        <Bar
          data={barData}
          options={{
            responsive: true,
            plugins: { legend: { display: false } },
            scales: {
              x: { ticks: { color: '#94a3b8' }, grid: { display: false } },
              y: { ticks: { color: '#94a3b8' }, grid: { color: 'rgba(148,163,184,0.1)' } },
            },
          }}
        />
      </Paper>
    </div>
  );
}
