import type { ButtonHTMLAttributes, ForwardedRef } from 'react';
import { forwardRef } from 'react';
import { cx } from './utils';

type ButtonVariant = 'primary' | 'secondary' | 'ghost' | 'destructive';
type ButtonSize = 'sm' | 'md' | 'lg' | 'icon';

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  size?: ButtonSize;
}

const variantStyles: Record<ButtonVariant, string> = {
  primary:
    'bg-accent text-slate-950 shadow-[0_14px_35px_-18px_rgba(78,162,255,0.9)] hover:bg-accent-strong',
  secondary: 'bg-surface-strong text-foreground hover:bg-slate-700/80',
  ghost: 'bg-transparent text-muted hover:bg-white/5 hover:text-foreground',
  destructive: 'bg-danger/15 text-danger hover:bg-danger/20 border border-danger/30',
};

const sizeStyles: Record<ButtonSize, string> = {
  sm: 'h-9 px-3 text-sm',
  md: 'h-11 px-4 text-sm',
  lg: 'h-12 px-5 text-base',
  icon: 'h-10 w-10 p-0',
};

export const Button = forwardRef(function Button(
  { className, variant = 'primary', size = 'md', type, ...props }: ButtonProps,
  ref: ForwardedRef<HTMLButtonElement>,
) {
  return (
    <button
      ref={ref}
      type={type ?? 'button'}
      className={cx(
        'inline-flex items-center justify-center gap-2 rounded-2xl font-semibold transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent/60 focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50',
        'border border-transparent',
        variantStyles[variant],
        sizeStyles[size],
        className,
      )}
      {...props}
    />
  );
});
