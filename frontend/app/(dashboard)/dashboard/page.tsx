'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { auth } from '@/lib/auth';
import { api, AccountResponse } from '@/lib/api';
import { Alert } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { EmptyState } from '@/components/ui/empty-state';
import { Skeleton } from '@/components/ui/skeleton';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';

export default function DashboardPage() {
  const router = useRouter();
  const [account, setAccount] = useState<AccountResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    let mounted = true;

    const loadAccount = async () => {
      const token = auth.getToken();

      if (!token) {
        router.replace('/login');
        return;
      }

      const response = await api.getAccount(token);

      if (!mounted) {
        return;
      }

      if (response.error) {
        setError(response.error);
        auth.logout();
        router.replace('/login');
        return;
      }

      if (response.data) {
        setAccount(response.data);
      }

      setLoading(false);
    };

    loadAccount();

    return () => {
      mounted = false;
    };
  }, [router]);

  const fullName = account ? `${account.firstName} ${account.lastName}` : '';

  const quickActions = [
    {
      title: 'Transfer',
      href: '/dashboard/transfer',
      description: 'Move money between accounts in a future release.',
      badge: 'Open',
    },
    {
      title: 'History',
      href: '/dashboard/history',
      description: 'Track payment activity and account changes here.',
      badge: 'Open',
    },
    {
      title: 'Profile',
      href: '/dashboard/profile',
      description: 'Update account details from the profile center.',
      badge: 'Open',
    },
    {
      title: 'Coupon',
      href: '/dashboard/coupon',
      description: 'Redeem offers and banking perks when it lands.',
      badge: 'Open',
    },
    {
      title: 'WebAuthn',
      href: '/dashboard/profile',
      description: 'Register a passkey from your secure profile center.',
      badge: 'New',
    },
  ] as const;

  const activityRows = [
    {
      event: 'Secure login session established',
      status: 'Complete',
      detail: 'JWT session checked from localStorage',
      time: 'Just now',
    },
    {
      event: 'Account summary loaded',
      status: 'Synced',
      detail: 'Balance and profile details are current',
      time: 'Moments ago',
    },
    {
      event: 'Transaction history slot reserved',
      status: 'Preview',
      detail: 'Future activity feed will appear here',
      time: 'Ready for phase 2',
    },
  ] as const;

  if (loading) {
    return (
      <div className="space-y-6 lg:space-y-8">
        <Card>
          <CardHeader className="gap-5 sm:flex-row sm:items-end sm:justify-between">
            <div className="space-y-3">
              <Skeleton className="h-4 w-36" />
              <Skeleton className="h-10 w-72 max-w-full" />
              <Skeleton className="h-5 w-96 max-w-full" />
            </div>
            <Skeleton className="h-10 w-32" />
          </CardHeader>
          <CardContent className="grid gap-4 sm:grid-cols-3">
            <Skeleton className="h-28 rounded-[1.25rem]" />
            <Skeleton className="h-28 rounded-[1.25rem]" />
            <Skeleton className="h-28 rounded-[1.25rem]" />
          </CardContent>
        </Card>

        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-5">
          {Array.from({ length: 5 }).map((_, index) => (
            <Skeleton key={index} className="h-32 rounded-[1.25rem]" />
          ))}
        </div>

        <Card>
          <CardHeader>
            <Skeleton className="h-4 w-40" />
            <Skeleton className="h-5 w-80 max-w-full" />
          </CardHeader>
          <CardContent className="space-y-3">
            <Skeleton className="h-14 w-full rounded-2xl" />
            <Skeleton className="h-14 w-full rounded-2xl" />
            <Skeleton className="h-14 w-full rounded-2xl" />
          </CardContent>
        </Card>
      </div>
    );
  }

  if (!account) {
    return (
      <Card>
        <CardContent className="p-6 sm:p-8">
          <EmptyState
            title="Unable to load your dashboard"
            description="Your session is active, but the account summary did not return. Try refreshing or signing in again."
          />
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-6 lg:space-y-8">
      {error ? (
        <Alert variant="destructive" title="Dashboard session issue">
          {error}
        </Alert>
      ) : null}

      <Card className="overflow-hidden border-accent/20 bg-[linear-gradient(135deg,rgba(14,22,41,0.98),rgba(11,17,32,0.9))]">
        <CardContent className="relative grid gap-6 p-6 sm:p-8 lg:grid-cols-[1.6fr_1fr] lg:gap-8 lg:p-10">
          <div className="absolute inset-0 -z-0 bg-[radial-gradient(circle_at_top_right,rgba(78,162,255,0.24),transparent_30%),radial-gradient(circle_at_bottom_left,rgba(57,217,138,0.12),transparent_28%),linear-gradient(180deg,rgba(255,255,255,0.02),transparent_35%)]" />
          <div className="relative z-10 space-y-6">
            <div className="flex flex-wrap items-center gap-3">
              <Badge variant="default">Private banking</Badge>
              <Badge variant="outline">Account {account.number}</Badge>
            </div>

            <div className="space-y-3">
              <p className="text-sm uppercase tracking-[0.28em] text-muted">Welcome back</p>
              <h1 className="text-3xl font-semibold tracking-tight text-foreground sm:text-4xl lg:text-5xl">
                {fullName}
              </h1>
              <p className="max-w-2xl text-sm leading-7 text-slate-200/80 sm:text-base">
                Your GoBank home is ready. Review your balance, jump into future actions, and keep everything aligned from one secure place.
              </p>
            </div>

            <div className="grid gap-3 sm:grid-cols-2">
              <div className="rounded-2xl border border-white/10 bg-white/[0.03] p-4">
                <p className="text-xs uppercase tracking-[0.14em] text-slate-300/70">Account Status</p>
                <p className="mt-2 text-sm font-semibold text-white">Active & Verified</p>
              </div>
              <div className="rounded-2xl border border-white/10 bg-white/[0.03] p-4">
                <p className="text-xs uppercase tracking-[0.14em] text-slate-300/70">Member Since</p>
                <p className="mt-2 text-sm font-semibold text-white">{new Date(account.createdAt).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}</p>
              </div>
              <div className="rounded-2xl border border-white/10 bg-white/[0.03] p-4">
                <p className="text-xs uppercase tracking-[0.14em] text-slate-300/70">Account Type</p>
                <p className="mt-2 text-sm font-semibold text-white">Premium Banking</p>
              </div>
              <div className="rounded-2xl border border-white/10 bg-white/[0.03] p-4">
                <p className="text-xs uppercase tracking-[0.14em] text-slate-300/70">Security</p>
                <p className="mt-2 text-sm font-semibold text-success">2FA Ready</p>
              </div>
            </div>

            <div className="space-y-2 rounded-2xl border border-accent/20 bg-[linear-gradient(180deg,rgba(78,162,255,0.12),rgba(78,162,255,0.04))] p-4">
              <p className="text-xs uppercase tracking-[0.14em] font-semibold text-accent">Why Choose GoBank?</p>
              <ul className="space-y-1 text-xs text-slate-200/80">
                <li className="flex items-start gap-2">
                  <span className="mt-1 text-accent">✓</span>
                  <span>Bank-grade security with passkey authentication</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="mt-1 text-accent">✓</span>
                  <span>Instant transfers between accounts</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="mt-1 text-accent">✓</span>
                  <span>Exclusive offers and rewards</span>
                </li>
              </ul>
            </div>

            <div className="flex flex-wrap gap-3">
              <Button variant="secondary" size="md" onClick={() => router.push('/dashboard/profile')}>
                View profile
              </Button>
              <Button variant="ghost" size="md" onClick={() => router.push('/dashboard/history')}>
                Explore activity
              </Button>
            </div>
          </div>

          <div className="relative z-10 grid gap-4 sm:grid-cols-2 lg:grid-cols-1">
            <Card className="border-accent/20 bg-[linear-gradient(180deg,rgba(78,162,255,0.16),rgba(255,255,255,0.04))] shadow-[0_26px_80px_-46px_rgba(78,162,255,0.5)]">
              <CardHeader className="pb-3">
                <CardDescription className="text-slate-200/80">Current balance</CardDescription>
                <CardTitle className="text-3xl font-bold tracking-tight text-white sm:text-4xl">${account.balance.toFixed(2)}</CardTitle>
              </CardHeader>
              <CardContent className="pt-0">
                <p className="text-sm leading-6 text-slate-100/80">Available funds for transfers, coupons, and everyday banking.</p>
              </CardContent>
            </Card>

            <Card className="border-accent/15 bg-[linear-gradient(180deg,rgba(17,26,46,0.96),rgba(12,19,35,0.96))] shadow-[0_22px_70px_-46px_rgba(14,165,233,0.45)]">
              <CardHeader className="pb-3">
                <CardDescription className="text-slate-200/80">Account summary</CardDescription>
                <CardTitle className="text-xl font-semibold text-white">{account.email}</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3 pt-0">
                <div className="flex items-center justify-between rounded-2xl border border-white/6 bg-white/[0.04] px-4 py-3 text-sm">
                  <span className="font-medium text-slate-200/75">Account number</span>
                  <span className="font-semibold tracking-[0.14em] text-white">{account.number}</span>
                </div>
                <div className="flex items-center justify-between rounded-2xl border border-white/6 bg-white/[0.03] px-4 py-3 text-sm">
                  <span className="font-medium text-slate-200/75">Profile status</span>
                  <Badge variant="success">Verified</Badge>
                </div>
              </CardContent>
            </Card>
          </div>
        </CardContent>
      </Card>

      <section className="space-y-3">
        <div className="flex items-end justify-between gap-3">
          <div>
            <p className="text-sm uppercase tracking-[0.28em] text-muted">Quick actions</p>
            <h2 className="mt-2 text-xl font-semibold text-foreground sm:text-2xl">Move into the next banking tasks</h2>
          </div>
        </div>

        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-5">
          {quickActions.map((action) => (
            <Card key={action.title} className="h-full border-white/8 bg-[linear-gradient(180deg,rgba(15,24,42,0.95),rgba(10,16,30,0.95))] transition-transform duration-200 hover:-translate-y-1 hover:border-accent/25 hover:bg-surface-strong/90">
              <CardContent className="flex h-full flex-col justify-between gap-4 p-5 sm:p-6">
                <div className="space-y-3">
                  <div className="flex items-center justify-between gap-3">
                    <div className="flex h-11 w-11 items-center justify-center rounded-2xl border border-accent/20 bg-[linear-gradient(180deg,rgba(78,162,255,0.22),rgba(78,162,255,0.08))] text-sm font-semibold text-accent">
                      {action.title.slice(0, 1)}
                    </div>
                    <Badge variant="outline">{action.badge}</Badge>
                  </div>
                  <div>
                    <h3 className="text-base font-semibold text-foreground">{action.title}</h3>
                    <p className="mt-2 text-sm leading-6 text-muted">{action.description}</p>
                  </div>
                </div>
                <Button variant="secondary" size="sm" className="w-full justify-between" onClick={() => router.push(action.href)}>
                  <span>{action.badge}</span>
                  <span aria-hidden="true">↗</span>
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      </section>

      <Card className="border-accent/12 bg-[linear-gradient(180deg,rgba(12,19,35,0.96),rgba(9,14,26,0.96))]">
        <CardHeader className="gap-3 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <p className="text-sm uppercase tracking-[0.28em] text-muted">Recent activity</p>
            <CardTitle className="mt-2 text-xl text-white sm:text-2xl">Preview of transaction history</CardTitle>
            <CardDescription>
              This section is ready for the transaction feed once transfers, offers, and history are wired up.
            </CardDescription>
          </div>
          <Badge variant="secondary">Placeholder feed</Badge>
        </CardHeader>
        <CardContent>
          <div className="overflow-hidden rounded-[1.25rem] border border-border/60 bg-white/[0.02]">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Activity</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Details</TableHead>
                  <TableHead className="text-right">Time</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {activityRows.map((row) => (
                  <TableRow key={row.event}>
                    <TableCell className="font-medium text-foreground">{row.event}</TableCell>
                    <TableCell>
                      <Badge variant={row.status === 'Complete' ? 'success' : 'secondary'}>{row.status}</Badge>
                    </TableCell>
                    <TableCell className="text-muted">{row.detail}</TableCell>
                    <TableCell className="text-right text-muted">{row.time}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>

          <div className="mt-5">
            <EmptyState
              title="Transaction history will populate here"
              description="Transfers, offers, and passkey activity can all flow into this area without redesigning the shell."
              icon={<span className="text-2xl">⟡</span>}
            />
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
