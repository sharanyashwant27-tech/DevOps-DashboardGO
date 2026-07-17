import { useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Alert,
  Button,
  MenuItem,
  Paper,
  Tab,
  Tabs,
  TextField,
  Typography,
} from '@mui/material';
import { k8sApi } from '../services/api';

export default function KubernetesPage() {
  const [namespace, setNamespace] = useState('default');
  const [tab, setTab] = useState(0);
  const [logs, setLogs] = useState('');
  const qc = useQueryClient();

  const nsQuery = useQuery({
    queryKey: ['k8s-ns'],
    queryFn: async () => (await k8sApi.namespaces()).data.data || [],
  });

  const dataQuery = useQuery({
    queryKey: ['k8s-data', tab, namespace],
    queryFn: async () => {
      if (tab === 0) return (await k8sApi.pods(namespace)).data.data;
      if (tab === 1) return (await k8sApi.deployments(namespace)).data.data;
      if (tab === 2) return (await k8sApi.services(namespace)).data.data;
      if (tab === 3) return (await k8sApi.nodes()).data.data;
      return (await k8sApi.events(namespace)).data.data;
    },
  });

  const restart = useMutation({
    mutationFn: (name: string) => k8sApi.restart(name, namespace),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['k8s-data'] }),
  });

  return (
    <div className="space-y-4">
      <Typography variant="h4" className="font-display">
        Kubernetes Cluster
      </Typography>
      {nsQuery.isError && <Alert severity="warning">Kubernetes API unavailable or disabled.</Alert>}

      <TextField
        select
        size="small"
        label="Namespace"
        value={namespace}
        onChange={(e) => setNamespace(e.target.value)}
        className="min-w-[200px]"
      >
        {(nsQuery.data as Array<{ metadata?: { name: string } }> | undefined)?.map((ns) => (
          <MenuItem key={ns.metadata?.name} value={ns.metadata?.name}>
            {ns.metadata?.name}
          </MenuItem>
        )) || <MenuItem value="default">default</MenuItem>}
      </TextField>

      <Tabs value={tab} onChange={(_, v) => setTab(v)}>
        <Tab label="Pods" />
        <Tab label="Deployments" />
        <Tab label="Services" />
        <Tab label="Nodes" />
        <Tab label="Events" />
      </Tabs>

      <Paper className="glass-panel p-4" elevation={0}>
        {tab === 1 && (
          <div className="mb-3 flex gap-2">
            <Button
              size="small"
              variant="outlined"
              onClick={() => {
                const items = dataQuery.data as Array<{ metadata?: { name: string } }>;
                if (items?.[0]?.metadata?.name) restart.mutate(items[0].metadata.name);
              }}
            >
              Restart first deployment
            </Button>
          </div>
        )}
        {tab === 0 && (
          <div className="mb-3">
            <Button
              size="small"
              onClick={async () => {
                const items = dataQuery.data as Array<{ metadata?: { name: string } }>;
                const pod = items?.[0]?.metadata?.name;
                if (!pod) return;
                const res = await k8sApi.podLogs(pod, namespace);
                setLogs((res.data.data as { log: string }).log || '');
              }}
            >
              Load first pod logs
            </Button>
          </div>
        )}
        <pre className="text-xs font-mono overflow-auto max-h-[60vh]">{JSON.stringify(dataQuery.data, null, 2)}</pre>
      </Paper>

      {logs && (
        <Paper className="glass-panel p-4" elevation={0}>
          <Typography variant="h6">Pod Logs</Typography>
          <pre className="text-xs font-mono overflow-auto max-h-80 whitespace-pre-wrap">{logs}</pre>
        </Paper>
      )}
    </div>
  );
}
