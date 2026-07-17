import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Button,
  Chip,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Typography,
} from '@mui/material';
import PageHero from '../components/common/PageHero';
import { deploymentApi } from '../services/api';
import type { Deployment } from '../types';

export default function DeploymentsPage() {
  const qc = useQueryClient();
  const query = useQuery({
    queryKey: ['deployments'],
    queryFn: async () => (await deploymentApi.list()).data.data || [],
  });
  const rollback = useMutation({
    mutationFn: (id: string) => deploymentApi.rollback(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['deployments'] }),
  });

  return (
    <div className="space-y-4">
      <PageHero
        eyebrow="Releases"
        title="Deployment History"
        subtitle="Track versions, environments, and rollbacks across applications"
        accent="#a3e635"
      />
      <Paper className="glass-panel overflow-auto" elevation={0}>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>Application</TableCell>
              <TableCell>Environment</TableCell>
              <TableCell>Version</TableCell>
              <TableCell>Commit</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Triggered By</TableCell>
              <TableCell>Time</TableCell>
              <TableCell align="right">Rollback</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {(query.data as Deployment[] | undefined)?.map((d) => (
              <TableRow key={d.id}>
                <TableCell>{d.application}</TableCell>
                <TableCell>{d.environment}</TableCell>
                <TableCell>{d.version}</TableCell>
                <TableCell className="font-mono text-xs">{d.git_commit}</TableCell>
                <TableCell>
                  <Chip size="small" label={d.status} />
                </TableCell>
                <TableCell>{d.triggered_by}</TableCell>
                <TableCell>{new Date(d.deployed_at).toLocaleString()}</TableCell>
                <TableCell align="right">
                  <Button
                    size="small"
                    disabled={!d.rollback_version}
                    onClick={() => rollback.mutate(d.id)}
                  >
                    Rollback
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
