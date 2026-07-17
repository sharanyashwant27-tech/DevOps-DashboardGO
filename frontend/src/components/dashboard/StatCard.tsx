import { Typography } from '@mui/material';
import type { ReactNode } from 'react';
import { Link } from 'react-router-dom';

interface Props {
  title: string;
  value: string | number;
  subtitle?: string;
  icon?: ReactNode;
  accent?: string;
  delay?: number;
  to?: string;
}

export default function StatCard({
  title,
  value,
  subtitle,
  icon,
  accent = '#22d3ee',
  delay = 0,
  to,
}: Props) {
  const content = (
    <>
      <div
        className="absolute inset-x-0 top-0 h-1"
        style={{ background: `linear-gradient(90deg, ${accent}, transparent 85%)` }}
      />
      <div
        className="absolute -right-8 -top-10 h-28 w-28 rounded-full blur-2xl opacity-40"
        style={{ background: accent }}
      />
      <div className="relative flex items-start justify-between gap-3">
        <div>
          <Typography variant="caption" color="text.secondary" className="uppercase tracking-wider">
            {title}
          </Typography>
          <Typography variant="h4" className="font-display mt-1" sx={{ color: accent }}>
            {value}
          </Typography>
          {subtitle && (
            <Typography variant="body2" color="text.secondary" className="mt-1">
              {subtitle}
            </Typography>
          )}
          {to && (
            <Typography variant="caption" sx={{ color: accent }} className="mt-2 inline-block font-medium">
              View details →
            </Typography>
          )}
        </div>
        <div
          className="rounded-2xl p-2.5 shadow-sm"
          style={{
            background: `linear-gradient(145deg, ${accent}33, ${accent}14)`,
            color: accent,
            border: `1px solid ${accent}44`,
          }}
        >
          {icon}
        </div>
      </div>
    </>
  );

  const className =
    'stat-card glass-panel relative overflow-hidden block no-underline text-inherit ' +
    (to ? 'cursor-pointer focus:outline-none focus-visible:ring-2 ' : '');

  const style = {
    animation: `rise 0.5s ease-out ${delay}ms both`,
    ['--card-accent' as string]: accent,
  };

  if (to) {
    return (
      <Link to={to} className={className} style={style}>
        {content}
      </Link>
    );
  }

  return (
    <div className={className} style={style}>
      {content}
    </div>
  );
}
