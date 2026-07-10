import type { HTMLAttributes } from 'react';
import { cx } from './utils';

export type CardVariant = 'default' | 'elevated' | 'featured' | 'gradient' | 'interactive';

interface CardProps extends HTMLAttributes<HTMLDivElement> {
  variant?: CardVariant;
}

const variantStyles: Record<CardVariant, string> = {
  default:
    'rounded-[1.5rem] border border-border/70 bg-surface/85 shadow-[0_24px_80px_-42px_rgba(2,6,23,0.95)] backdrop-blur transition-all duration-300',
  elevated:
    'rounded-[1.5rem] border border-border/50 bg-surface-strong/85 shadow-[0_32px_100px_-50px_rgba(78,162,255,0.4)] backdrop-blur transition-all duration-300',
  featured:
    'rounded-[1.5rem] border border-featured/40 bg-gradient-to-br from-featured/8 via-surface-strong/85 to-surface/85 shadow-[0_32px_100px_-50px_rgba(244,63,94,0.3)] backdrop-blur transition-all duration-300',
  gradient:
    'rounded-[1.5rem] border border-accent/30 bg-gradient-to-br from-primary/10 via-surface-strong/85 to-accent-secondary/5 shadow-[0_32px_100px_-50px_rgba(78,162,255,0.35)] backdrop-blur transition-all duration-300',
  interactive:
    'rounded-[1.5rem] border border-border/70 bg-surface/85 shadow-[0_24px_80px_-42px_rgba(2,6,23,0.95)] backdrop-blur transition-all duration-300 hover:border-accent/50 hover:-translate-y-1 hover:shadow-[0_32px_100px_-50px_rgba(78,162,255,0.4)] cursor-pointer',
};

export function Card({ className, variant = 'default', ...props }: CardProps) {
  return (
    <div
      className={cx(
        variantStyles[variant],
        className,
      )}
      {...props}
    />
  );
}

export function CardHeader({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div 
      className={cx(
        'flex flex-col gap-2 p-6 sm:p-7 border-b border-border/40 bg-gradient-to-r from-primary/5 to-transparent',
        className
      )} 
      {...props} 
    />
  );
}

export function CardTitle({ className, ...props }: HTMLAttributes<HTMLHeadingElement>) {
  return <h3 className={cx('text-lg font-semibold tracking-tight text-foreground animate-entrance', className)} {...props} />;
}

export function CardDescription({ className, ...props }: HTMLAttributes<HTMLParagraphElement>) {
  return <p className={cx('text-sm text-muted leading-6', className)} {...props} />;
}

export function CardContent({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return <div className={cx('px-6 pb-6 sm:px-7 sm:pb-7 animate-entrance-up', className)} {...props} />;
}

export function CardFooter({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div 
      className={cx(
        'flex items-center gap-3 border-t border-border/40 px-6 py-5 sm:px-7 bg-surface-soft/50 rounded-b-[1.25rem] transition-colors duration-300',
        className
      )} 
      {...props} 
    />
  );
}
