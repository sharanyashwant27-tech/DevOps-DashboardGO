import { Typography } from '@mui/material';
import type { ReactNode } from 'react';

interface Props {
  eyebrow: string;
  title: string;
  subtitle?: string;
  accent?: string;
  action?: ReactNode;
}

export default function PageHero({ eyebrow, title, subtitle, accent = '#22d3ee', action }: Props) {
  return (
    <div className="page-hero flex items-start justify-between gap-4 flex-wrap">
      <div className="relative z-[1]">
        <Typography variant="overline" sx={{ color: accent, letterSpacing: 2 }}>
          {eyebrow}
        </Typography>
        <Typography variant="h4" className="font-display">
          {title}
        </Typography>
        {subtitle && (
          <Typography color="text.secondary" className="max-w-2xl">
            {subtitle}
          </Typography>
        )}
      </div>
      {action && <div className="relative z-[1]">{action}</div>}
    </div>
  );
}
