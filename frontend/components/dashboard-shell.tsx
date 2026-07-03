'use client';

import type { ReactNode } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { auth } from '@/lib/auth';

interface DashboardShellProps {
  children: ReactNode;
}

export function DashboardShell({ children }: DashboardShellProps) {
  const router = useRouter();
  const user = auth.getUser();

  const handleLogout = () => {
    auth.logout();
    router.replace('/login');
  };

  return (
    <div className="gobank-shell">
      <header className="sticky top-0 z-40 border-b border-border/70 bg-background/80 backdrop-blur-xl">
        <div className="mx-auto flex w-full max-w-7xl items-center justify-between gap-4 px-4 py-4 sm:px-6 lg:px-8">
          <div className="flex items-center gap-4">
            <div className="flex h-12 w-12 items-center justify-center rounded-2xl border border-accent/20 bg-accent/12 text-lg font-bold text-accent shadow-[0_14px_40px_-22px_rgba(78,162,255,0.75)]">
              G
            </div>
            <div>
              <p className="text-lg font-semibold tracking-tight text-foreground">GoBank</p>
              <div className="mt-1 flex flex-wrap items-center gap-2">
                <Badge variant="secondary">Secure session</Badge>
                <Badge variant="outline">{user?.number ? `Acct ${user.number}` : 'Protected banking'}</Badge>
              </div>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <Link href="/dashboard" className="rounded-full border border-border/70 bg-white/[0.03] px-4 py-2 text-sm font-semibold text-slate-200 transition hover:border-accent/25 hover:text-white">
              Home
            </Link>
            <Link href="/dashboard/transfer" className="rounded-full border border-border/70 bg-white/[0.03] px-4 py-2 text-sm font-semibold text-slate-200 transition hover:border-accent/25 hover:text-white">
              Transfer
            </Link>
            <Link href="/dashboard/history" className="rounded-full border border-border/70 bg-white/[0.03] px-4 py-2 text-sm font-semibold text-slate-200 transition hover:border-accent/25 hover:text-white">
              History
            </Link>
            <Link href="/dashboard/profile" className="rounded-full border border-border/70 bg-white/[0.03] px-4 py-2 text-sm font-semibold text-slate-200 transition hover:border-accent/25 hover:text-white">
              Profile
            </Link>
            <Link href="/dashboard/coupon" className="rounded-full border border-border/70 bg-white/[0.03] px-4 py-2 text-sm font-semibold text-slate-200 transition hover:border-accent/25 hover:text-white">
              Coupon
            </Link>
            <Button variant="destructive" size="sm" onClick={handleLogout} className="shrink-0">
              Logout
            </Button>
          </div>
        </div>
      </header>

      <main className="mx-auto w-full max-w-7xl flex-1 px-4 py-6 sm:px-6 lg:px-8 lg:py-8">{children}</main>
    </div>
  );
}
