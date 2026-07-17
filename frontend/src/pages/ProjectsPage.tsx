import { useQuery } from '@tanstack/react-query';
import {
  Chip,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  TextField,
  Typography,
} from '@mui/material';
import { useState } from 'react';
import { Link } from 'react-router-dom';
import PageHero from '../components/common/PageHero';
import { projectApi } from '../services/api';
import type { Project } from '../types';

export default function ProjectsPage() {
  const [search, setSearch] = useState('');
  const query = useQuery({
    queryKey: ['projects', search],
    queryFn: async () => (await projectApi.list(search)).data.data || [],
  });

  return (
    <div className="space-y-4">
      <PageHero
        eyebrow="Portfolio"
        title="Projects"
        subtitle="Organizations and applications managed in Command Center"
        accent="#38bdf8"
      />

      <TextField
        size="small"
        label="Search projects"
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        className="max-w-sm"
      />

      <Paper className="glass-panel overflow-auto" elevation={0}>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Slug</TableCell>
              <TableCell>Environment</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Repository</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {(query.data as Project[] | undefined)?.map((p) => (
              <TableRow key={p.id} hover>
                <TableCell>
                  <Link to={`/deployments`} className="text-teal-400 hover:underline font-medium">
                    {p.name}
                  </Link>
                </TableCell>
                <TableCell className="font-mono text-xs">{p.slug}</TableCell>
                <TableCell>
                  <Chip size="small" label={p.environment} variant="outlined" />
                </TableCell>
                <TableCell>
                  <Chip size="small" label={p.status} color={p.status === 'active' ? 'success' : 'default'} />
                </TableCell>
                <TableCell>
                  {p.repository_url ? (
                    <a href={p.repository_url} target="_blank" rel="noreferrer" className="text-teal-400 hover:underline text-sm">
                      {p.repository_url}
                    </a>
                  ) : (
                    '-'
                  )}
                </TableCell>
              </TableRow>
            ))}
            {!query.isLoading && !(query.data as Project[] | undefined)?.length && (
              <TableRow>
                <TableCell colSpan={5}>
                  <Typography color="text.secondary">No projects found.</Typography>
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </Paper>
    </div>
  );
}
