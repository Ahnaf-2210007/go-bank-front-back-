import type { ReactNode } from 'react';
import { Button } from './button';
import { cx } from './utils';

export type ModalVariant = 'default' | 'featured' | 'alert';

export interface ModalProps {
  open: boolean;
  title: ReactNode;
  description?: ReactNode;
  children?: ReactNode;
  onClose: () => void;
  variant?: ModalVariant;
  closeButton?: boolean;
}

const variantStyles: Record<ModalVariant, string> = {
  default: 'border-border/70 bg-surface shadow-[0_40px_120px_-48px_rgba(2,6,23,1)]',
  featured: 'border-featured/40 bg-gradient-to-br from-featured/8 via-surface to-surface shadow-[0_40px_120px_-48px_rgba(244,63,94,0.3)]',
  alert: 'border-danger/40 bg-gradient-to-br from-danger/8 via-surface to-surface shadow-[0_40px_120px_-48px_rgba(239,68,68,0.3)]',
};

export function Modal({ open, title, description, children, onClose, variant = 'default', closeButton = true }: ModalProps) {
  if (!open) {
    return null;
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center px-4 py-6 animate-fade-in">
      <button
        type="button"
        className="absolute inset-0 bg-slate-950/75 backdrop-blur-sm transition-opacity duration-300 animate-fade-in"
        aria-label="Close modal overlay"
        onClick={onClose}
      />
      <div className={cx('relative z-10 w-full max-w-lg rounded-[1.75rem] border p-6 sm:p-8 transition-all duration-300 animate-fade-scale-in', variantStyles[variant])}>
        <div className="space-y-2 animate-entrance">
          <h2 className="text-xl font-semibold text-foreground animate-entrance-up">{title}</h2>
          {description ? <p className="text-sm leading-6 text-muted animate-entrance-up">{description}</p> : null}
        </div>
        <div className="mt-5 text-sm text-slate-200/90 animate-entrance-up">{children}</div>
        {closeButton && (
          <div className="mt-6 flex justify-end gap-3 animate-entrance-up">
            <Button variant="secondary" onClick={onClose}>
              Close
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}
