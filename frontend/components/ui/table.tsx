import type { HTMLAttributes, TableHTMLAttributes } from 'react';
import { cx } from './utils';

export type TableVariant = 'default' | 'striped' | 'bordered' | 'minimal';

interface TableProps extends TableHTMLAttributes<HTMLTableElement> {
  variant?: TableVariant;
}

const variantStyles: Record<TableVariant, string> = {
  default: 'w-full border-separate border-spacing-0',
  striped: 'w-full border-separate border-spacing-0 [&_tbody_tr:nth-child(even)]:bg-white/[0.02]',
  bordered: 'w-full border-collapse border border-border/50',
  minimal: 'w-full border-separate border-spacing-0',
};

export function Table({ className, variant = 'default', ...props }: TableProps) {
  return <table className={cx(variantStyles[variant], className)} {...props} />;
}

export function TableHeader({ className, ...props }: HTMLAttributes<HTMLTableSectionElement>) {
  return <thead className={cx('text-left text-xs uppercase tracking-[0.18em] text-muted/80 bg-gradient-to-r from-primary/8 to-transparent', className)} {...props} />;
}

export function TableBody({ className, ...props }: HTMLAttributes<HTMLTableSectionElement>) {
  return <tbody className={cx('text-sm text-slate-100 animate-stagger', className)} {...props} />;
}

export function TableRow({ className, ...props }: HTMLAttributes<HTMLTableRowElement>) {
  return (
    <tr 
      className={cx(
        'transition-all duration-300 hover:bg-primary/5 animate-entrance-up',
        className
      )} 
      {...props} 
    />
  );
}

export function TableHead({ className, ...props }: HTMLAttributes<HTMLTableCellElement>) {
  return (
    <th 
      className={cx(
        'border-b border-border/50 px-4 py-3 font-semibold text-primary-light text-left transition-colors duration-300',
        className
      )} 
      {...props} 
    />
  );
}

export function TableCell({ className, ...props }: HTMLAttributes<HTMLTableCellElement>) {
  return (
    <td 
      className={cx(
        'border-b border-border/30 px-4 py-4 align-middle transition-colors duration-300',
        className
      )} 
      {...props} 
    />
  );
}
