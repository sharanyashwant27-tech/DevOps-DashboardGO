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
  accent = '#0ea5a4',
  delay = 0,
  to,
}: Props) {
  const content = (
    <>
      <div
        className="absolute -right-6 -top-6 h-24 w-24 rounded-full opacity-20"
        style={{ background: accent }}
      />
      <div className="flex items-start justify-between gap-3">
        <div>
          <Typography variant="caption" color="text.secondary" className="uppercase tracking-wider">
            {title}
          </Typography>
          <Typography variant="h4" className="font-display mt-1">
            {value}
          </Typography>
          {subtitle && (
            <Typography variant="body2" color="text.secondary" className="mt-1">
              {subtitle}
            </Typography>
          )}
          {to && (
            <Typography variant="caption" sx={{ color: accent }} className="mt-2 inline-block opacity-80">
              View details →
            </Typography>
          )}
        </div>
        <div className="rounded-xl p-2" style={{ background: `${accent}22`, color: accent }}>
          {icon}
        </div>
      </div>
    </>
  );

  const className =
    'stat-card glass-panel relative overflow-hidden block no-underline text-inherit ' +
    (to ? 'cursor-pointer hover:ring-1 focus:outline-none focus-visible:ring-2 ' : '');

  if (to) {
    return (
      <Link
        to={to}
        className={className}
        style={{ animation: `rise 0.5s ease-out ${delay}ms both` }}
      >
        {content}
      </Link>
    );
  }

  return (
    <div className={className} style={{ animation: `rise 0.5s ease-out ${delay}ms both` }}>
      {content}
    </div>
  );
}
