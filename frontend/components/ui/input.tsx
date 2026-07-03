import type { ForwardedRef, InputHTMLAttributes } from 'react';
import { forwardRef } from 'react';
import { cx } from './utils';

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {}

export const Input = forwardRef(function Input(
  { className, ...props }: InputProps,
  ref: ForwardedRef<HTMLInputElement>,
) {
  return (
    <input
      ref={ref}
      className={cx(
        'h-11 w-full rounded-2xl border border-border/70 bg-surface-strong/85 px-4 text-sm text-foreground shadow-inner outline-none transition placeholder:text-slate-500 focus:border-accent/70 focus:ring-2 focus:ring-accent/20 disabled:cursor-not-allowed disabled:opacity-60',
        className,
      )}
      {...props}
    />
  );
});
