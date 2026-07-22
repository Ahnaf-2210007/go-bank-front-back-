import type { HTMLAttributes } from 'react';
import { cx } from './utils';

type SkeletonVariant = 'default' | 'text' | 'circle' | 'avatar' | 'card';

interface SkeletonProps extends HTMLAttributes<HTMLDivElement> {
  variant?: SkeletonVariant;
}

const variantStyles: Record<SkeletonVariant, string> = {
  default: 'rounded-2xl h-10 bg-gradient-to-r from-primary/15 via-primary/8 to-primary/15 animate-shimmer',
  text: 'rounded-lg h-4 bg-gradient-to-r from-primary/15 via-primary/8 to-primary/15 animate-shimmer',
  circle: 'rounded-full bg-gradient-to-r from-primary/15 via-primary/8 to-primary/15 animate-shimmer',
  avatar: 'rounded-full w-12 h-12 bg-gradient-to-r from-primary/15 via-primary/8 to-primary/15 animate-shimmer',
  card: 'rounded-2xl bg-gradient-to-br from-primary/10 via-surface/50 to-primary/10 animate-pulse-smooth',
};

export function Skeleton({ className, variant = 'default', ...props }: SkeletonProps) {
  return (
    <div 
      className={cx(variantStyles[variant], className)} 
      {...props} 
    />
  );
}

export function SkeletonGroup({ count = 3, variant = 'text', gap = 'gap-3' }: { count?: number; variant?: SkeletonVariant; gap?: string }) {
  return (
    <div className={`space-y-${gap}`}>
      {Array.from({ length: count }).map((_, i) => (
        <Skeleton key={i} variant={variant} className="w-full" />
      ))}
    </div>
  );
}
