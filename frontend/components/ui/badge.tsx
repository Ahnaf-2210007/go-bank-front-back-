import type { HTMLAttributes } from 'react';
import { cx } from './utils';

type BadgeVariant = 'default' | 'secondary' | 'success' | 'warning' | 'destructive' | 'outline';

export interface BadgeProps extends HTMLAttributes<HTMLSpanElement> {
  variant?: BadgeVariant;
}

const badgeVariants: Record<BadgeVariant, string> = {
  default: 'bg-accent/15 text-accent border border-accent/25',
  secondary: 'bg-white/[0.06] text-slate-200 border border-white/10',
  success: 'bg-success/15 text-success border border-success/25',
  warning: 'bg-warning/15 text-warning border border-warning/25',
  destructive: 'bg-danger/15 text-danger border border-danger/25',
  outline: 'bg-transparent text-slate-200 border border-border/70',
};

export function Badge({ className, variant = 'default', ...props }: BadgeProps) {
  return (
    <span
      className={cx(
        'inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold uppercase tracking-[0.18em]',
        badgeVariants[variant],
        className,
      )}
      {...props}
    />
  );
}
