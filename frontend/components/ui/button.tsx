import type { ButtonHTMLAttributes, ForwardedRef } from 'react';
import { forwardRef } from 'react';
import { cx } from './utils';

type ButtonVariant = 'primary' | 'secondary' | 'ghost' | 'destructive' | 'gradient' | 'success' | 'warning' | 'premium';
type ButtonSize = 'sm' | 'md' | 'lg' | 'icon';

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  size?: ButtonSize;
  isLoading?: boolean;
}

const variantStyles: Record<ButtonVariant, string> = {
  primary:
    'bg-accent text-slate-950 shadow-[0_14px_35px_-18px_rgba(78,162,255,0.9)] hover:bg-accent-strong hover:shadow-[0_18px_45px_-22px_rgba(78,162,255,0.95)] active:scale-95',
  secondary: 'bg-surface-strong text-foreground hover:bg-slate-700/80 hover:shadow-md active:scale-95',
  ghost: 'bg-transparent text-muted hover:bg-white/5 hover:text-foreground hover:shadow-sm active:bg-white/10',
  destructive: 'bg-danger/15 text-danger hover:bg-danger/25 border border-danger/30 hover:border-danger/50 hover:shadow-[0_8px_16px_-4px_rgba(239,68,68,0.3)] active:scale-95',
  gradient: 'bg-gradient-to-br from-primary via-primary-dark to-accent-secondary text-white shadow-[0_14px_35px_-18px_rgba(78,162,255,0.7)] hover:shadow-[0_20px_45px_-22px_rgba(78,162,255,0.8)] active:scale-95',
  success: 'bg-success text-white shadow-[0_14px_35px_-18px_rgba(16,185,129,0.6)] hover:bg-success-dark hover:shadow-[0_18px_45px_-22px_rgba(16,185,129,0.7)] active:scale-95',
  warning: 'bg-warning text-slate-950 shadow-[0_14px_35px_-18px_rgba(245,158,11,0.5)] hover:bg-warning-dark hover:shadow-[0_18px_45px_-22px_rgba(245,158,11,0.6)] active:scale-95',
  premium: 'bg-gradient-to-br from-premium via-accent-secondary to-featured text-white shadow-[0_14px_35px_-18px_rgba(139,92,246,0.6)] hover:shadow-[0_20px_45px_-22px_rgba(139,92,246,0.7)] active:scale-95',
};

const sizeStyles: Record<ButtonSize, string> = {
  sm: 'h-9 px-3 text-sm gap-1.5',
  md: 'h-11 px-4 text-sm gap-2',
  lg: 'h-12 px-6 text-base gap-2.5',
  icon: 'h-10 w-10 p-0',
};

export const Button = forwardRef(function Button(
  { className, variant = 'primary', size = 'md', type, isLoading, disabled, ...props }: ButtonProps,
  ref: ForwardedRef<HTMLButtonElement>,
) {
  return (
    <button
      ref={ref}
      type={type ?? 'button'}
      disabled={isLoading || disabled}
      className={cx(
        'inline-flex items-center justify-center rounded-2xl font-semibold transition-all duration-250 ease-out focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50',
        'border border-transparent',
        variantStyles[variant],
        sizeStyles[size],
        'active:transition-none',
        isLoading && 'opacity-80 cursor-wait',
        className,
      )}
      {...props}
    >
      {isLoading && (
        <svg className="animate-spin-smooth -ml-1 mr-2 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
        </svg>
      )}
      {props.children}
    </button>
  );
});
