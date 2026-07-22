import type { ForwardedRef, InputHTMLAttributes } from 'react';
import { forwardRef } from 'react';
import { cx } from './utils';

export type InputVariant = 'default' | 'success' | 'error' | 'warning';

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  variant?: InputVariant;
  hasError?: boolean;
  isSuccess?: boolean;
}

const variantStyles: Record<InputVariant, string> = {
  default: 'border-border/70 focus:border-accent/70 focus:ring-accent/20',
  success: 'border-success/50 focus:border-success/80 focus:ring-success/20 bg-success/[0.03]',
  error: 'border-danger/50 focus:border-danger/80 focus:ring-danger/20 bg-danger/[0.03]',
  warning: 'border-warning/50 focus:border-warning/80 focus:ring-warning/20 bg-warning/[0.03]',
};

export const Input = forwardRef(function Input(
  { className, variant = 'default', hasError, isSuccess, ...props }: InputProps,
  ref: ForwardedRef<HTMLInputElement>,
) {
  const computedVariant = hasError ? 'error' : isSuccess ? 'success' : variant;
  
  return (
    <input
      ref={ref}
      className={cx(
        'h-11 w-full rounded-2xl border bg-surface-strong/85 px-4 text-sm text-foreground shadow-inner outline-none transition-all duration-250 placeholder:text-slate-500 focus:ring-2 focus:shadow-[0_0_12px_rgba(78,162,255,0.15)] disabled:cursor-not-allowed disabled:opacity-60 animate-entrance-up',
        'hover:border-border/100 focus:bg-surface-strong/95',
        variantStyles[computedVariant],
        className,
      )}
      {...props}
    />
  );
});
