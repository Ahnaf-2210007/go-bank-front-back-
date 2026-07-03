import type { HTMLAttributes } from 'react';
import { cx } from './utils';

export function Card({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cx(
        'rounded-[1.5rem] border border-border/70 bg-surface/85 shadow-[0_24px_80px_-42px_rgba(2,6,23,0.95)] backdrop-blur',
        className,
      )}
      {...props}
    />
  );
}

export function CardHeader({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return <div className={cx('flex flex-col gap-2 p-6 sm:p-7', className)} {...props} />;
}

export function CardTitle({ className, ...props }: HTMLAttributes<HTMLHeadingElement>) {
  return <h3 className={cx('text-lg font-semibold tracking-tight text-foreground', className)} {...props} />;
}

export function CardDescription({ className, ...props }: HTMLAttributes<HTMLParagraphElement>) {
  return <p className={cx('text-sm text-muted leading-6', className)} {...props} />;
}

export function CardContent({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return <div className={cx('px-6 pb-6 sm:px-7 sm:pb-7', className)} {...props} />;
}

export function CardFooter({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return <div className={cx('flex items-center gap-3 border-t border-border/60 px-6 py-5 sm:px-7', className)} {...props} />;
}
