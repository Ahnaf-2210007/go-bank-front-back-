import type { ReactNode } from 'react';
import { Button } from './button';
import { cx } from './utils';

export interface ModalProps {
  open: boolean;
  title: ReactNode;
  description?: ReactNode;
  children?: ReactNode;
  onClose: () => void;
}

export function Modal({ open, title, description, children, onClose }: ModalProps) {
  if (!open) {
    return null;
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center px-4 py-6">
      <button
        type="button"
        className="absolute inset-0 bg-slate-950/75 backdrop-blur-sm"
        aria-label="Close modal overlay"
        onClick={onClose}
      />
      <div className={cx('relative z-10 w-full max-w-lg rounded-[1.75rem] border border-border/70 bg-surface p-6 shadow-[0_40px_120px_-48px_rgba(2,6,23,1)]')}>
        <div className="space-y-2">
          <h2 className="text-xl font-semibold text-foreground">{title}</h2>
          {description ? <p className="text-sm leading-6 text-muted">{description}</p> : null}
        </div>
        <div className="mt-5 text-sm text-slate-200/90">{children}</div>
        <div className="mt-6 flex justify-end">
          <Button variant="secondary" onClick={onClose}>Close</Button>
        </div>
      </div>
    </div>
  );
}
