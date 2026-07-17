import { FormEvent, useState } from 'react';
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
import { incidentApi } from '../services/api';
import type { Incident } from '../types';

export default function IncidentsPage() {
  const qc = useQueryClient();
  const [params, setParams] = useSearchParams();
  const [title, setTitle] = useState('');
  const [priority, setPriority] = useState('high');
  const [statusFilter, setStatusFilter] = useState(params.get('status') || '');

  const query = useQuery({
    queryKey: ['incidents', statusFilter],
    queryFn: async () =>
      (await incidentApi.list(statusFilter ? { status: statusFilter } : undefined)).data.data || [],
  });

  const create = useMutation({
    mutationFn: () => incidentApi.create({ title, priority, description: title }),
    onSuccess: () => {
      setTitle('');
      qc.invalidateQueries({ queryKey: ['incidents'] });
    },
  });

  const resolve = useMutation({
    mutationFn: (id: string) => incidentApi.update(id, { status: 'resolved', resolution: 'Mitigated' }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['incidents'] }),
  });

  function onSubmit(e: FormEvent) {
    e.preventDefault();
    create.mutate();
  }

  return (
    <div className="space-y-4">
      <Typography variant="h4" className="font-display">
        Incident Management
      </Typography>

      <TextField
        select
        size="small"
        label="Status filter"
        value={statusFilter}
        onChange={(e) => {
          const v = e.target.value;
          setStatusFilter(v);
          if (!v) params.delete('status');
          else params.set('status', v);
          setParams(params);
        }}
        className="min-w-[180px]"
      >
        <MenuItem value="">All</MenuItem>
        <MenuItem value="open">Open</MenuItem>
        <MenuItem value="investigating">Investigating</MenuItem>
        <MenuItem value="mitigated">Mitigated</MenuItem>
        <MenuItem value="resolved">Resolved</MenuItem>
        <MenuItem value="closed">Closed</MenuItem>
      </TextField>

      <Paper className="glass-panel p-4" elevation={0} component="form" onSubmit={onSubmit}>
        <Typography variant="h6" className="mb-3">
          Create Incident
        </Typography>
        <div className="flex flex-col md:flex-row gap-3">
          <TextField label="Title" value={title} onChange={(e) => setTitle(e.target.value)} required fullWidth />
          <TextField select label="Priority" value={priority} onChange={(e) => setPriority(e.target.value)} className="min-w-[160px]">
            <MenuItem value="critical">Critical</MenuItem>
            <MenuItem value="high">High</MenuItem>
            <MenuItem value="medium">Medium</MenuItem>
            <MenuItem value="low">Low</MenuItem>
          </TextField>
          <Button type="submit" variant="contained">
            Create
          </Button>
        </div>
      </Paper>

      <Paper className="glass-panel overflow-auto" elevation={0}>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>Title</TableCell>
              <TableCell>Priority</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>SLA</TableCell>
              <TableCell align="right">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {(query.data as Incident[] | undefined)?.map((i) => (
              <TableRow key={i.id}>
                <TableCell>{i.title}</TableCell>
                <TableCell>
                  <Chip size="small" label={i.priority} color={i.priority === 'critical' ? 'error' : 'default'} />
                </TableCell>
                <TableCell>{i.status}</TableCell>
                <TableCell>{i.sla_deadline ? new Date(i.sla_deadline).toLocaleString() : '-'}</TableCell>
                <TableCell align="right">
                  <Button size="small" onClick={() => resolve.mutate(i.id)}>
                    Resolve
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
