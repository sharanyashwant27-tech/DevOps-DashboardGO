import { useMemo, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Alert,
  Button,
  Chip,
  MenuItem,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  TextField,
  Typography,
} from '@mui/material';
import { Link, useSearchParams } from 'react-router-dom';
import PlayCircleIcon from '@mui/icons-material/PlayCircle';
import ErrorIcon from '@mui/icons-material/Error';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import ViewListIcon from '@mui/icons-material/ViewList';
import QueueIcon from '@mui/icons-material/Queue';
import { jenkinsApi } from '../services/api';
import StatCard from '../components/dashboard/StatCard';

type JenkinsJob = { name: string; color: string; url: string };

function matchesFilter(color: string, filter: string) {
  const c = (color || '').toLowerCase();
  if (!filter) return true;
  if (filter === 'running') return c.includes('anime');
  if (filter === 'failed') return c.startsWith('red') && !c.includes('anime');
  if (filter === 'success') return c.startsWith('blue') && !c.includes('anime');
  if (filter === 'queue') return false;
  return true;
}

export default function JenkinsPage() {
  const [params, setParams] = useSearchParams();
  const [search, setSearch] = useState('');
  const [consoleLog, setConsoleLog] = useState('');
  const [selectedJob, setSelectedJob] = useState(params.get('job') || '');
  const filter = params.get('filter') || '';
  const qc = useQueryClient();

  const jobsQuery = useQuery({
    queryKey: ['jenkins-jobs', search],
    queryFn: async () => (await jenkinsApi.jobs(search)).data.data || [],
  });

  const jobs = useMemo(() => {
    const items = (jobsQuery.data as JenkinsJob[] | undefined) || [];
    if (filter === 'queue') return items;
    return items.filter((j) => matchesFilter(j.color, filter));
  }, [jobsQuery.data, filter]);

  const queueQuery = useQuery({
    queryKey: ['jenkins-queue'],
    queryFn: async () => (await jenkinsApi.queue()).data.data || [],
  });
  const statsQuery = useQuery({
    queryKey: ['jenkins-stats'],
    queryFn: async () => (await jenkinsApi.stats()).data.data,
  });

  const buildsQuery = useQuery({
    queryKey: ['jenkins-builds', selectedJob],
    enabled: !!selectedJob,
    queryFn: async () => (await jenkinsApi.builds(selectedJob)).data.data || [],
  });

  const trigger = useMutation({
    mutationFn: (job: string) => jenkinsApi.trigger(job),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['jenkins-jobs'] });
      qc.invalidateQueries({ queryKey: ['jenkins-stats'] });
      qc.invalidateQueries({ queryKey: ['jenkins-builds'] });
    },
  });

  function setFilter(next: string) {
    if (!next) params.delete('filter');
    else params.set('filter', next);
    setParams(params);
  }

  function openJob(name: string) {
    setSelectedJob(name);
    params.set('job', name);
    setParams(params);
  }

  const stats = statsQuery.data as {
    total_jobs?: number;
    running_jobs?: number;
    failed_jobs?: number;
    successful_jobs?: number;
    mode?: string;
  } | undefined;

  const showQueue = filter === 'queue';

  return (
    <div className="space-y-4">
      <div className="flex items-start justify-between gap-3 flex-wrap">
        <div>
          <Typography variant="h4" className="font-display">
            Jenkins Pipelines
          </Typography>
          <Typography color="text.secondary">
            Click a status card to filter jobs — open a job row for builds and console
          </Typography>
        </div>
        <Button component={Link} to="/settings" size="small" variant="outlined">
          Integration settings
        </Button>
      </div>

      {filter && (
        <Chip
          label={`Filter: ${filter}`}
          color="primary"
          variant="outlined"
          onDelete={() => setFilter('')}
        />
      )}
      {jobsQuery.isError && <Alert severity="warning">Jenkins unavailable or not configured.</Alert>}
      {!jobsQuery.isError && stats?.mode === 'demo' && (
        <Alert severity="info">
          Jenkins demo mode is active (sample jobs). Set <code>jenkins.url</code> /{' '}
          <code>DCC_JENKINS_URL</code> to connect a real Jenkins server.
        </Alert>
      )}

      <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-5 gap-4">
        <StatCard
          title="Total Jobs"
          value={stats?.total_jobs ?? 0}
          icon={<ViewListIcon />}
          to="/jenkins"
          delay={0}
        />
        <StatCard
          title="Running"
          value={stats?.running_jobs ?? 0}
          icon={<PlayCircleIcon />}
          accent="#3b82f6"
          to="/jenkins?filter=running"
          delay={40}
        />
        <StatCard
          title="Failed"
          value={stats?.failed_jobs ?? 0}
          icon={<ErrorIcon />}
          accent="#ef4444"
          to="/jenkins?filter=failed"
          delay={80}
        />
        <StatCard
          title="Successful"
          value={stats?.successful_jobs ?? 0}
          icon={<CheckCircleIcon />}
          accent="#22c55e"
          to="/jenkins?filter=success"
          delay={120}
        />
        <StatCard
          title="Build Queue"
          value={(queueQuery.data as unknown[] | undefined)?.length ?? 0}
          icon={<QueueIcon />}
          accent="#f59e0b"
          to="/jenkins?filter=queue"
          delay={160}
        />
      </div>

      {showQueue ? (
        <Paper className="glass-panel p-4" elevation={0} id="queue">
          <Typography variant="h6" className="mb-2">
            Build Queue ({(queueQuery.data as unknown[])?.length || 0})
          </Typography>
          <pre className="text-xs font-mono overflow-auto max-h-80">
            {JSON.stringify(queueQuery.data, null, 2)}
          </pre>
        </Paper>
      ) : (
        <>
          <div className="flex flex-wrap gap-3">
            <TextField
              size="small"
              label="Search jobs"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="max-w-sm"
            />
            <TextField
              select
              size="small"
              label="Status filter"
              value={filter === 'queue' ? '' : filter}
              onChange={(e) => setFilter(e.target.value)}
              className="min-w-[160px]"
            >
              <MenuItem value="">All</MenuItem>
              <MenuItem value="running">Running</MenuItem>
              <MenuItem value="failed">Failed</MenuItem>
              <MenuItem value="success">Successful</MenuItem>
            </TextField>
          </div>

          <Paper className="glass-panel overflow-auto" elevation={0}>
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>Job</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell align="right">Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {jobs.map((job) => (
                  <TableRow
                    key={job.name}
                    hover
                    selected={selectedJob === job.name}
                    className="cursor-pointer"
                    onClick={() => openJob(job.name)}
                  >
                    <TableCell>
                      <button
                        type="button"
                        className="text-teal-400 hover:underline font-medium bg-transparent border-0 p-0 cursor-pointer"
                        onClick={(e) => {
                          e.stopPropagation();
                          openJob(job.name);
                        }}
                      >
                        {job.name}
                      </button>
                    </TableCell>
                    <TableCell>
                      <Chip size="small" label={job.color} />
                    </TableCell>
                    <TableCell align="right" className="space-x-2" onClick={(e) => e.stopPropagation()}>
                      <Button size="small" onClick={() => trigger.mutate(job.name)}>
                        Trigger
                      </Button>
                      <Button
                        size="small"
                        onClick={async () => {
                          openJob(job.name);
                          const builds = (await jenkinsApi.builds(job.name)).data.data as Array<{ number: number }>;
                          if (builds?.[0]) {
                            const log = await jenkinsApi.console(job.name, builds[0].number);
                            setConsoleLog((log.data.data as { log: string }).log || '');
                          }
                        }}
                      >
                        Console
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
                {!jobsQuery.isLoading && jobs.length === 0 && (
                  <TableRow>
                    <TableCell colSpan={3}>
                      <Typography color="text.secondary">
                        {jobsQuery.isError
                          ? 'Configure Jenkins to view jobs, or clear the filter.'
                          : 'No jobs match this filter.'}
                      </Typography>
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </Paper>
        </>
      )}

      {selectedJob && !showQueue && (
        <Paper className="glass-panel p-4" elevation={0}>
          <div className="flex items-center justify-between mb-2 gap-2 flex-wrap">
            <Typography variant="h6">Build history — {selectedJob}</Typography>
            <Button size="small" onClick={() => { setSelectedJob(''); params.delete('job'); setParams(params); }}>
              Close
            </Button>
          </div>
          <pre className="text-xs font-mono overflow-auto max-h-60">
            {JSON.stringify(buildsQuery.data, null, 2)}
          </pre>
        </Paper>
      )}

      {consoleLog && (
        <Paper className="glass-panel p-4" elevation={0}>
          <Typography variant="h6" className="mb-2">
            Console Output
          </Typography>
          <pre className="text-xs font-mono overflow-auto max-h-80 whitespace-pre-wrap">{consoleLog}</pre>
        </Paper>
      )}
    </div>
  );
}
