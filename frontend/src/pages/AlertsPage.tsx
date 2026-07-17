import { useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
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
import { useSearchParams } from 'react-router-dom';
import PageHero from '../components/common/PageHero';
import { alertApi } from '../services/api';
import type { Alert } from '../types';

export default function AlertsPage() {
  const [params, setParams] = useSearchParams();
  const [severity, setSeverity] = useState(params.get('severity') || '');
  const [status, setStatus] = useState(params.get('status') || '');
  const [search, setSearch] = useState(params.get('search') || '');
  const qc = useQueryClient();

  function updateSeverity(v: string) {
    setSeverity(v);
    if (!v) params.delete('severity');
    else params.set('severity', v);
    setParams(params);
  }
  function updateStatus(v: string) {
    setStatus(v);
    if (!v) params.delete('status');
    else params.set('status', v);
    setParams(params);
  }

  const query = useQuery({
    queryKey: ['alerts', severity, status, search],
    queryFn: async () =>
      (
        await alertApi.list({
          ...(severity ? { severity } : {}),
          ...(status ? { status } : {}),
          ...(search ? { search } : {}),
        })
      ).data.data || [],
  });

  const ack = useMutation({
    mutationFn: (id: string) => alertApi.acknowledge(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['alerts'] }),
  });
  const resolve = useMutation({
    mutationFn: (id: string) => alertApi.resolve(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['alerts'] }),
  });
  const mute = useMutation({
    mutationFn: (id: string) => alertApi.mute(id, 60),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['alerts'] }),
  });

  return (
    <div className="space-y-4">
      <PageHero
        eyebrow="Signals"
        title="Alert Dashboard"
        subtitle="Acknowledge, mute, and resolve alerts across severity levels"
        accent="#f59e0b"
      />
      <div className="flex flex-wrap gap-3">
        <TextField size="small" label="Search" value={search} onChange={(e) => setSearch(e.target.value)} />
        <TextField select size="small" label="Severity" value={severity} onChange={(e) => updateSeverity(e.target.value)} className="min-w-[140px]">
          <MenuItem value="">All</MenuItem>
          <MenuItem value="critical">Critical</MenuItem>
          <MenuItem value="high">High</MenuItem>
          <MenuItem value="medium">Medium</MenuItem>
          <MenuItem value="low">Low</MenuItem>
        </TextField>
        <TextField select size="small" label="Status" value={status} onChange={(e) => updateStatus(e.target.value)} className="min-w-[140px]">
          <MenuItem value="">All</MenuItem>
          <MenuItem value="open">Open</MenuItem>
          <MenuItem value="acknowledged">Acknowledged</MenuItem>
          <MenuItem value="resolved">Resolved</MenuItem>
          <MenuItem value="muted">Muted</MenuItem>
        </TextField>
      </div>

      <Paper className="glass-panel overflow-auto" elevation={0}>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>Title</TableCell>
              <TableCell>Severity</TableCell>
              <TableCell>Source</TableCell>
              <TableCell>Status</TableCell>
              <TableCell align="right">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {(query.data as Alert[] | undefined)?.map((a) => (
              <TableRow key={a.id}>
                <TableCell>{a.title}</TableCell>
                <TableCell>
                  <Chip size="small" label={a.severity} color={a.severity === 'critical' ? 'error' : 'warning'} />
                </TableCell>
                <TableCell>{a.source}</TableCell>
                <TableCell>{a.status}</TableCell>
                <TableCell align="right" className="space-x-1">
                  <Button size="small" onClick={() => ack.mutate(a.id)}>
                    Ack
                  </Button>
                  <Button size="small" onClick={() => resolve.mutate(a.id)}>
                    Resolve
                  </Button>
                  <Button size="small" onClick={() => mute.mutate(a.id)}>
                    Mute
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </Paper>
    </div>
  );
}
