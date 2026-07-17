import { Chip } from '@mui/material';

export default function JobStatusChip({ color }: { color: string }) {
  const label = color.includes('anime') ? 'running' : color.startsWith('blue') ? 'success' : color;
  return <Chip size="small" label={label} color={label === 'running' ? 'info' : 'default'} />;
}
