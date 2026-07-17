import { Chip } from '@mui/material';

export default function JobStatusChip({ color }: { color: string }) {
  const label = color.includes('anime') ? 'running' : color.startsWith('blue') ? 'success' : color;
  const styles: Record<string, { bg: string; fg: string; border: string }> = {
    running: { bg: 'rgba(56,189,248,0.18)', fg: '#7dd3fc', border: 'rgba(56,189,248,0.45)' },
    success: { bg: 'rgba(52,211,153,0.18)', fg: '#6ee7b7', border: 'rgba(52,211,153,0.45)' },
    failure: { bg: 'rgba(251,113,133,0.18)', fg: '#fda4af', border: 'rgba(251,113,133,0.45)' },
    failed: { bg: 'rgba(251,113,133,0.18)', fg: '#fda4af', border: 'rgba(251,113,133,0.45)' },
    aborted: { bg: 'rgba(148,163,184,0.16)', fg: '#cbd5e1', border: 'rgba(148,163,184,0.4)' },
    unstable: { bg: 'rgba(251,191,36,0.18)', fg: '#fcd34d', border: 'rgba(251,191,36,0.45)' },
  };
  const s = styles[label] || { bg: 'rgba(34,211,238,0.14)', fg: '#67e8f9', border: 'rgba(34,211,238,0.35)' };
  return (
    <Chip
      size="small"
      label={label}
      sx={{
        bgcolor: s.bg,
        color: s.fg,
        border: `1px solid ${s.border}`,
        fontWeight: 600,
        textTransform: 'capitalize',
      }}
    />
  );
}
