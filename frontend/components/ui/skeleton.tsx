import type { HTMLAttributes } from 'react';
import { cx } from './utils';

export function Skeleton({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return <div className={cx('animate-pulse rounded-2xl bg-white/[0.08]', className)} {...props} />;
}
