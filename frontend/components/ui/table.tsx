import type { HTMLAttributes, TableHTMLAttributes } from 'react';
import { cx } from './utils';

export function Table({ className, ...props }: TableHTMLAttributes<HTMLTableElement>) {
  return <table className={cx('w-full border-separate border-spacing-0', className)} {...props} />;
}

export function TableHeader({ className, ...props }: HTMLAttributes<HTMLTableSectionElement>) {
  return <thead className={cx('text-left text-xs uppercase tracking-[0.18em] text-muted', className)} {...props} />;
}

export function TableBody({ className, ...props }: HTMLAttributes<HTMLTableSectionElement>) {
  return <tbody className={cx('text-sm text-slate-100', className)} {...props} />;
}

export function TableRow({ className, ...props }: HTMLAttributes<HTMLTableRowElement>) {
  return <tr className={cx('transition-colors hover:bg-white/[0.03]', className)} {...props} />;
}

export function TableHead({ className, ...props }: HTMLAttributes<HTMLTableCellElement>) {
  return <th className={cx('border-b border-border/60 px-4 py-3 font-medium', className)} {...props} />;
}

export function TableCell({ className, ...props }: HTMLAttributes<HTMLTableCellElement>) {
  return <td className={cx('border-b border-border/50 px-4 py-4 align-middle', className)} {...props} />;
}
