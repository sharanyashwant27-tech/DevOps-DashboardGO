import { Chip } from '@mui/material';

export default function AlertBadge({ severity }: { severity: string }) {
  const color =
    severity === 'critical' ? 'error' : severity === 'high' ? 'warning' : 'default';
  return <Chip size="small" label={severity} color={color as 'error' | 'warning' | 'default'} />;
}
