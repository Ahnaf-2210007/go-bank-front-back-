import type { ReactNode } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

interface AuthPageProps {
  title: string;
  description: string;
  children: ReactNode;
  footer?: ReactNode;
  eyebrow?: ReactNode;
}

export function AuthPage({ title, description, children, footer, eyebrow }: AuthPageProps) {
  return (
    <div className="relative min-h-screen overflow-hidden px-4 py-8 sm:px-6 lg:px-8">
      <div className="absolute inset-0 -z-0 bg-[radial-gradient(circle_at_top_left,rgba(78,162,255,0.18),transparent_26%),radial-gradient(circle_at_bottom_right,rgba(57,217,138,0.08),transparent_24%)]" />
      <div className="mx-auto flex min-h-[calc(100vh-4rem)] w-full max-w-5xl items-center justify-center">
        <div className="grid w-full gap-6 lg:grid-cols-[1fr_1.1fr] lg:gap-8">
          <div className="gobank-panel-strong flex flex-col justify-between p-6 sm:p-8 lg:p-10">
            <div className="space-y-5">
              <div className="flex h-12 w-12 items-center justify-center rounded-2xl border border-accent/20 bg-accent/12 text-lg font-bold text-accent shadow-[0_14px_40px_-22px_rgba(78,162,255,0.75)]">
                G
              </div>
              <div className="space-y-3">
                <p className="text-sm uppercase tracking-[0.28em] text-muted">GoBank secure access</p>
                <h1 className="text-4xl font-semibold tracking-tight text-white sm:text-5xl">{title}</h1>
                <p className="max-w-md text-base leading-7 text-slate-200/75">{description}</p>
              </div>
              {eyebrow ? <div>{eyebrow}</div> : null}
            </div>

            <div className="mt-8 grid gap-3 sm:grid-cols-2">
              <div className="rounded-2xl border border-white/8 bg-white/[0.03] p-4">
                <p className="text-xs uppercase tracking-[0.22em] text-muted">JWT auth</p>
                <p className="mt-2 text-sm text-slate-100/80">localStorage session, server-verified on protected routes</p>
              </div>
              <div className="rounded-2xl border border-white/8 bg-white/[0.03] p-4">
                <p className="text-xs uppercase tracking-[0.22em] text-muted">Banking shell</p>
                <p className="mt-2 text-sm text-slate-100/80">consistent spacing, focus states, and clean dark surfaces</p>
              </div>
            </div>
          </div>

          <Card className="border-white/8 bg-[linear-gradient(180deg,rgba(12,19,35,0.97),rgba(9,14,26,0.97))]">
            <CardHeader className="space-y-2 pb-0">
              <CardTitle className="text-2xl text-white sm:text-3xl">{title}</CardTitle>
              <CardDescription className="text-slate-200/75">{description}</CardDescription>
            </CardHeader>
            <CardContent className="pt-6">{children}</CardContent>
            {footer ? <div className="border-t border-border/60 px-6 py-5 sm:px-8">{footer}</div> : null}
          </Card>
        </div>
      </div>
    </div>
  );
}
