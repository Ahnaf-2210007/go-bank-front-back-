import type { ReactNode } from 'react';
import { Button } from './button';
import { cx } from './utils';

export interface EmptyStateProps {
  title: ReactNode;
  description: ReactNode;
  actionLabel?: ReactNode;
  onAction?: () => void;
  icon?: ReactNode;
  className?: string;
}

export function EmptyState({ title, description, actionLabel, onAction, icon, className }: EmptyStateProps) {
  return (
    <div className={cx('rounded-[1.5rem] border border-dashed border-border/70 bg-surface/55 p-6 text-center sm:p-8', className)}>
      {icon ? <div className="mb-4 flex justify-center text-accent">{icon}</div> : null}
      <h3 className="text-lg font-semibold text-foreground">{title}</h3>
      <p className="mx-auto mt-2 max-w-xl text-sm leading-6 text-muted">{description}</p>
      {actionLabel && onAction ? (
        <div className="mt-5">
          <Button variant="secondary" onClick={onAction}>
            {actionLabel}
          </Button>
        </div>
      ) : null}
    </div>
  );
}
