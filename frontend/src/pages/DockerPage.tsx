import { useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Alert,
  Button,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  TextField,
  Typography,
} from '@mui/material';
import { dockerApi } from '../services/api';

export default function DockerPage() {
  const [search, setSearch] = useState('');
  const [logs, setLogs] = useState('');
  const qc = useQueryClient();

  const containers = useQuery({
    queryKey: ['docker-containers', search],
    queryFn: async () => (await dockerApi.containers(search)).data.data || [],
  });
  const images = useQuery({
    queryKey: ['docker-images'],
    queryFn: async () => (await dockerApi.images()).data.data || [],
  });

  const action = useMutation({
    mutationFn: async ({ id, op }: { id: string; op: 'start' | 'stop' | 'restart' | 'remove' }) => {
      if (op === 'start') return dockerApi.start(id);
      if (op === 'stop') return dockerApi.stop(id);
      if (op === 'restart') return dockerApi.restart(id);
      return dockerApi.remove(id);
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ['docker-containers'] }),
  });

  return (
    <div className="space-y-4">
      <Typography variant="h4" className="font-display">
        Docker Monitoring
      </Typography>
      {containers.isError && <Alert severity="warning">Docker daemon unavailable.</Alert>}

      <TextField size="small" label="Search containers" value={search} onChange={(e) => setSearch(e.target.value)} />

      <Paper className="glass-panel overflow-auto" elevation={0}>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Image</TableCell>
              <TableCell>State</TableCell>
              <TableCell>Status</TableCell>
              <TableCell align="right">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {(containers.data as Array<{ Id: string; Names: string[]; Image: string; State: string; Status: string }> | undefined)?.map(
              (c) => (
                <TableRow key={c.Id}>
                  <TableCell>{c.Names?.[0]?.replace(/^\//, '')}</TableCell>
                  <TableCell>{c.Image}</TableCell>
                  <TableCell>{c.State}</TableCell>
                  <TableCell>{c.Status}</TableCell>
                  <TableCell align="right" className="space-x-1">
                    <Button size="small" onClick={() => action.mutate({ id: c.Id, op: 'start' })}>
                      Start
                    </Button>
                    <Button size="small" onClick={() => action.mutate({ id: c.Id, op: 'stop' })}>
                      Stop
                    </Button>
                    <Button size="small" onClick={() => action.mutate({ id: c.Id, op: 'restart' })}>
                      Restart
                    </Button>
                    <Button
                      size="small"
                      onClick={async () => {
                        const res = await dockerApi.logs(c.Id);
                        setLogs((res.data.data as { log: string }).log || '');
                      }}
                    >
                      Logs
                    </Button>
                    <Button size="small" color="error" onClick={() => action.mutate({ id: c.Id, op: 'remove' })}>
                      Delete
                    </Button>
                  </TableCell>
                </TableRow>
              ),
            )}
          </TableBody>
        </Table>
      </Paper>

      <Paper className="glass-panel p-4" elevation={0}>
        <Typography variant="h6">Images ({(images.data as unknown[])?.length || 0})</Typography>
        <pre className="text-xs font-mono overflow-auto max-h-48">{JSON.stringify(images.data, null, 2)}</pre>
      </Paper>

      {logs && (
        <Paper className="glass-panel p-4" elevation={0}>
          <Typography variant="h6">Container Logs</Typography>
          <pre className="text-xs font-mono overflow-auto max-h-80 whitespace-pre-wrap">{logs}</pre>
        </Paper>
      )}
    </div>
  );
}
