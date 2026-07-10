import type { HTMLAttributes, ReactNode } from 'react';
import { cx } from './utils';

type AlertVariant = 'default' | 'success' | 'warning' | 'destructive' | 'info';

export interface AlertProps extends HTMLAttributes<HTMLDivElement> {
  variant?: AlertVariant;
  title?: ReactNode;
  icon?: ReactNode;
  isDismissible?: boolean;
  onDismiss?: () => void;
}

const alertStyles: Record<AlertVariant, string> = {
  default: 'border-accent/40 bg-gradient-to-r from-accent/15 to-accent/5 text-slate-100 shadow-[0_4px_12px_-2px_rgba(78,162,255,0.15)]',
  success: 'border-success/40 bg-gradient-to-r from-success/15 to-success/5 text-slate-100 shadow-[0_4px_12px_-2px_rgba(16,185,129,0.15)]',
  warning: 'border-warning/40 bg-gradient-to-r from-warning/15 to-warning/5 text-slate-100 shadow-[0_4px_12px_-2px_rgba(245,158,11,0.15)]',
  destructive: 'border-danger/40 bg-gradient-to-r from-danger/15 to-danger/5 text-slate-100 shadow-[0_4px_12px_-2px_rgba(239,68,68,0.15)]',
  info: 'border-primary/40 bg-gradient-to-r from-primary/15 to-primary/5 text-slate-100 shadow-[0_4px_12px_-2px_rgba(78,162,255,0.15)]',
};

const iconColors: Record<AlertVariant, string> = {
  default: 'text-accent',
  success: 'text-success',
  warning: 'text-warning',
  destructive: 'text-danger',
  info: 'text-primary',
};

export function Alert({ className, variant = 'default', title, icon, isDismissible, onDismiss, children, ...props }: AlertProps) {
  return (
    <div className={cx('rounded-3xl border px-4 py-3 sm:px-5 transition-all duration-300 animate-entrance', alertStyles[variant], className)} {...props}>
      <div className="flex items-start gap-3">
        {icon && <div className={cx('mt-0.5 text-lg flex-shrink-0', iconColors[variant])}>{icon}</div>}
        <div className="flex-1">
          {title ? <div className="text-sm font-semibold animate-entrance-up">{title}</div> : null}
          {children ? <div className={cx('text-sm leading-6', title ? 'mt-1 text-slate-200/90' : 'text-slate-200/90', 'animate-entrance-up')}>{children}</div> : null}
        </div>
        {isDismissible && onDismiss && (
          <button
            type="button"
            onClick={onDismiss}
            className="flex-shrink-0 mt-0.5 text-slate-400 hover:text-slate-200 transition-colors duration-200 text-lg leading-none"
            aria-label="Dismiss alert"
          >
            ×
          </button>
        )}
      </div>
    </div>
  );
}
