import { useQuery } from '@tanstack/react-query';
import { Paper, Typography } from '@mui/material';
import { serverApi } from '../services/api';

export default function ServersPage() {
  const list = useQuery({
    queryKey: ['servers'],
    queryFn: async () => (await serverApi.list()).data.data || [],
  });
  const local = useQuery({
    queryKey: ['server-local'],
    queryFn: async () => (await serverApi.local()).data.data,
  });

  return (
    <div className="space-y-4">
      <Typography variant="h4" className="font-display">
        Server Monitoring
      </Typography>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <Paper className="glass-panel p-4" elevation={0}>
          <Typography variant="h6" className="mb-2">
            Registered Servers
          </Typography>
          <pre className="text-xs font-mono overflow-auto max-h-[60vh]">{JSON.stringify(list.data, null, 2)}</pre>
        </Paper>
        <Paper className="glass-panel p-4" elevation={0}>
          <Typography variant="h6" className="mb-2">
            Local Host Details
          </Typography>
          <pre className="text-xs font-mono overflow-auto max-h-[60vh]">{JSON.stringify(local.data, null, 2)}</pre>
        </Paper>
      </div>
    </div>
  );
}
