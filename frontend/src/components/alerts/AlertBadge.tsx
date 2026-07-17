import { Chip } from '@mui/material';

export default function AlertBadge({ severity }: { severity: string }) {
  const styles: Record<string, { bg: string; fg: string; border: string }> = {
    critical: { bg: 'rgba(244,63,94,0.18)', fg: '#fb7185', border: 'rgba(244,63,94,0.45)' },
    high: { bg: 'rgba(251,191,36,0.18)', fg: '#fbbf24', border: 'rgba(251,191,36,0.45)' },
    medium: { bg: 'rgba(56,189,248,0.16)', fg: '#38bdf8', border: 'rgba(56,189,248,0.4)' },
    low: { bg: 'rgba(52,211,153,0.16)', fg: '#34d399', border: 'rgba(52,211,153,0.4)' },
  };
  const s = styles[severity] || styles.medium;
  return (
    <Chip
      size="small"
      label={severity}
      sx={{
        bgcolor: s.bg,
        color: s.fg,
        border: `1px solid ${s.border}`,
        fontWeight: 700,
        textTransform: 'capitalize',
      }}
    />
  );
}
