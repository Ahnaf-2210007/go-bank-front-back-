import type { HTMLAttributes } from 'react';
import { cx } from './utils';

type BadgeVariant = 'default' | 'secondary' | 'success' | 'warning' | 'destructive' | 'outline' | 'premium' | 'featured' | 'status-active' | 'status-pending';

export interface BadgeProps extends HTMLAttributes<HTMLSpanElement> {
  variant?: BadgeVariant;
  isAnimated?: boolean;
}

const badgeVariants: Record<BadgeVariant, string> = {
  default: 'bg-accent/15 text-accent border border-accent/25 shadow-[0_4px_12px_-2px_rgba(78,162,255,0.2)] transition-all duration-300',
  secondary: 'bg-white/[0.06] text-slate-200 border border-white/10 shadow-sm transition-all duration-300',
  success: 'bg-success/15 text-success border border-success/30 shadow-[0_4px_12px_-2px_rgba(16,185,129,0.2)] transition-all duration-300',
  warning: 'bg-warning/15 text-warning border border-warning/30 shadow-[0_4px_12px_-2px_rgba(245,158,11,0.2)] transition-all duration-300',
  destructive: 'bg-danger/15 text-danger border border-danger/30 shadow-[0_4px_12px_-2px_rgba(239,68,68,0.2)] transition-all duration-300',
  outline: 'bg-transparent text-slate-200 border border-border/70 transition-all duration-300 hover:border-border/100 hover:shadow-sm',
  premium: 'bg-gradient-to-r from-premium/20 to-accent-secondary/10 text-premium border border-premium/40 shadow-[0_4px_12px_-2px_rgba(139,92,246,0.3)] transition-all duration-300',
  featured: 'bg-gradient-to-r from-featured/20 to-warning/10 text-featured border border-featured/40 shadow-[0_4px_12px_-2px_rgba(244,63,94,0.3)] transition-all duration-300',
  'status-active': 'bg-success/10 text-success border border-success/50 shadow-[0_4px_12px_-2px_rgba(16,185,129,0.3)] animate-pulse-smooth',
  'status-pending': 'bg-warning/10 text-warning border border-warning/50 shadow-[0_4px_12px_-2px_rgba(245,158,11,0.3)] animate-pulse-smooth',
};

export function Badge({ className, variant = 'default', isAnimated, ...props }: BadgeProps) {
  return (
    <span
      className={cx(
        'inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-semibold uppercase tracking-[0.18em] animate-entrance',
        badgeVariants[variant],
        isAnimated && 'animate-pulse-smooth',
        className,
      )}
      {...props}
    />
  );
}
