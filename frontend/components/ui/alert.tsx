import type { HTMLAttributes, ReactNode } from 'react';
import { cx } from './utils';

type AlertVariant = 'default' | 'success' | 'warning' | 'destructive';

export interface AlertProps extends HTMLAttributes<HTMLDivElement> {
  variant?: AlertVariant;
  title?: ReactNode;
}

const alertStyles: Record<AlertVariant, string> = {
  default: 'border-accent/25 bg-accent/10 text-slate-100',
  success: 'border-success/25 bg-success/10 text-slate-100',
  warning: 'border-warning/25 bg-warning/10 text-slate-100',
  destructive: 'border-danger/25 bg-danger/10 text-slate-100',
};

export function Alert({ className, variant = 'default', title, children, ...props }: AlertProps) {
  return (
    <div className={cx('rounded-3xl border px-4 py-3 sm:px-5', alertStyles[variant], className)} {...props}>
      {title ? <div className="text-sm font-semibold">{title}</div> : null}
      {children ? <div className={cx('text-sm leading-6', title ? 'mt-1 text-slate-200/90' : 'text-slate-200/90')}>{children}</div> : null}
    </div>
  );
}
