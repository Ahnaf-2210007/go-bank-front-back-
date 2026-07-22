'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { auth } from '@/lib/auth';
import { api, AccountResponse } from '@/lib/api';
import { ActivityEvent } from '@/lib/types';
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
  const [activity, setActivity] = useState<ActivityEvent[]>([]);
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

      const [accountResponse, activityResponse] = await Promise.all([
        api.getAccount(token),
        api.getActivity(token),
      ]);

      if (!mounted) {
        return;
      }

      if (accountResponse.error) {
        setError(accountResponse.error);
        auth.logout();
        router.replace('/login');
        return;
      }

      if (accountResponse.data) {
        setAccount(accountResponse.data);
      }

      if (activityResponse.data) {
        setActivity(activityResponse.data);
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

  const getActivityTitle = (event: ActivityEvent) => {
    switch (event.type) {
      case 'transfer_sent':
        return event.title;
      case 'transfer_received':
        return event.title;
      case 'account_created':
        return 'Account Created';
      case 'passkey_registered':
        return 'Passkey Registered';
      case 'profile_update':
        return 'Profile Updated';
      default:
        return event.title;
    }
  };

  const getStatusBadgeVariant = (status: string) => {
    switch (status.toLowerCase()) {
      case 'completed':
        return 'success';
      case 'synced':
        return 'secondary';
      case 'pending':
        return 'secondary';
      default:
        return 'secondary';
    }
  };

  const formatTime = (timestamp: string) => {
    const now = new Date();
    const date = new Date(timestamp);
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
  };

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
    <div className="space-y-6 lg:space-y-8 animate-stagger">
      {error ? (
        <Alert variant="destructive" title="Dashboard session issue">
          {error}
        </Alert>
      ) : null}

      <Card variant="gradient" className="overflow-hidden">
        <CardContent className="relative grid gap-6 p-6 sm:p-8 lg:grid-cols-[1.6fr_1fr] lg:gap-8 lg:p-10">
          <div className="absolute inset-0 -z-0 bg-[radial-gradient(circle_at_top_right,rgba(78,162,255,0.24),transparent_30%),radial-gradient(circle_at_bottom_left,rgba(16,185,129,0.12),transparent_28%),linear-gradient(180deg,rgba(255,255,255,0.02),transparent_35%)]" />
          <div className="relative z-10 space-y-6 animate-entrance">
            <div className="flex flex-wrap items-center gap-3 animate-entrance-up">
              <Badge variant="premium">Private banking</Badge>
              <Badge variant="outline">Account {account.number}</Badge>
            </div>

            <div className="space-y-3 animate-entrance-up">
              <p className="text-sm uppercase tracking-[0.28em] text-muted animate-entrance">Welcome back</p>
              <h1 className="text-3xl font-semibold tracking-tight text-foreground sm:text-4xl lg:text-5xl animate-entrance-up">
                {fullName}
              </h1>
              <p className="max-w-2xl text-sm leading-7 text-slate-200/80 sm:text-base animate-entrance-up">
                Your GoBank home is ready. Review your balance, jump into future actions, and keep everything aligned from one secure place.
              </p>
            </div>

            <div className="grid gap-3 sm:grid-cols-2 animate-stagger">
              <div className="rounded-2xl border border-success/20 bg-success/[0.05] p-4 transition-all duration-300 hover:border-success/40 hover:shadow-[0_8px_16px_-4px_rgba(16,185,129,0.2)]">
                <p className="text-xs uppercase tracking-[0.14em] text-slate-300/70">Account Status</p>
                <p className="mt-2 text-sm font-semibold text-success-light">Active & Verified</p>
              </div>
              <div className="rounded-2xl border border-primary/20 bg-primary/[0.05] p-4 transition-all duration-300 hover:border-primary/40 hover:shadow-[0_8px_16px_-4px_rgba(78,162,255,0.2)]">
                <p className="text-xs uppercase tracking-[0.14em] text-slate-300/70">Member Since</p>
                <p className="mt-2 text-sm font-semibold text-primary-light">{new Date(account.createdAt).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}</p>
              </div>
              <div className="rounded-2xl border border-premium/20 bg-premium/[0.05] p-4 transition-all duration-300 hover:border-premium/40 hover:shadow-[0_8px_16px_-4px_rgba(139,92,246,0.2)]">
                <p className="text-xs uppercase tracking-[0.14em] text-slate-300/70">Account Type</p>
                <p className="mt-2 text-sm font-semibold text-premium-light">Premium Banking</p>
              </div>
              <div className="rounded-2xl border border-featured/20 bg-featured/[0.05] p-4 transition-all duration-300 hover:border-featured/40 hover:shadow-[0_8px_16px_-4px_rgba(244,63,94,0.2)]">
                <p className="text-xs uppercase tracking-[0.14em] text-slate-300/70">Security</p>
                <p className="mt-2 text-sm font-semibold text-success-light flex items-center gap-1.5"><span className="animate-pulse-smooth">●</span>2FA Ready</p>
              </div>
            </div>

            <div className="space-y-2 rounded-2xl border border-primary/30 bg-[linear-gradient(180deg,rgba(78,162,255,0.12),rgba(78,162,255,0.04))] p-4 transition-all duration-300 hover:border-primary/50 hover:shadow-[0_8px_20px_-4px_rgba(78,162,255,0.25)]">
              <p className="text-xs uppercase tracking-[0.14em] font-semibold text-primary-light animate-entrance">Why Choose GoBank?</p>
              <ul className="space-y-1 text-xs text-slate-200/80 animate-stagger">
                <li className="flex items-start gap-2">
                  <span className="mt-1 text-success animate-entrance">✓</span>
                  <span>Bank-grade security with passkey authentication</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="mt-1 text-success animate-entrance">✓</span>
                  <span>Instant transfers between accounts</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="mt-1 text-success animate-entrance">✓</span>
                  <span>Exclusive offers and rewards</span>
                </li>
              </ul>
            </div>

            <div className="flex flex-wrap gap-3 animate-entrance-up">
              <Button variant="gradient" size="md" onClick={() => router.push('/dashboard/profile')}>
                View profile
              </Button>
              <Button variant="secondary" size="md" onClick={() => router.push('/dashboard/history')}>
                Explore activity
              </Button>
            </div>
          </div>

          <div className="relative z-10 grid gap-4 sm:grid-cols-2 lg:grid-cols-1 animate-stagger">
            <Card variant="featured" className="border-primary/20 bg-gradient-to-br from-primary/15 via-primary-dark/8 to-primary/5 animate-hover-lift-strong">
              <CardHeader className="pb-3">
                <CardDescription className="text-slate-200/80">Current balance</CardDescription>
                <CardTitle className="text-3xl font-bold tracking-tight text-primary-light sm:text-4xl animate-entrance-up">${account.balance.toFixed(2)}</CardTitle>
              </CardHeader>
              <CardContent className="pt-0">
                <p className="text-sm leading-6 text-slate-100/80">Available funds for transfers, coupons, and everyday banking.</p>
              </CardContent>
            </Card>

            <Card variant="elevated" className="animate-hover-lift">
              <CardHeader className="pb-3 border-b border-border/40 bg-gradient-to-r from-accent-secondary/10 to-transparent">
                <CardDescription className="text-slate-200/80">Account summary</CardDescription>
                <CardTitle className="text-xl font-semibold text-white animate-entrance">{account.email}</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3 pt-3 animate-stagger">
                <div className="flex items-center justify-between rounded-2xl border border-primary/20 bg-primary/[0.05] px-4 py-3 text-sm transition-all duration-300 hover:border-primary/40 hover:bg-primary/[0.08]">
                  <span className="font-medium text-slate-200/75">Account number</span>
                  <span className="font-semibold tracking-[0.14em] text-primary-light">{account.number}</span>
                </div>
                <div className="flex items-center justify-between rounded-2xl border border-success/20 bg-success/[0.05] px-4 py-3 text-sm transition-all duration-300 hover:border-success/40 hover:bg-success/[0.08]">
                  <span className="font-medium text-slate-200/75">Profile status</span>
                  <Badge variant="status-active">Verified</Badge>
                </div>
              </CardContent>
            </Card>
          </div>
        </CardContent>
      </Card>

      <section className="space-y-3 animate-entrance-up">
        <div className="flex items-end justify-between gap-3">
          <div>
            <p className="text-sm uppercase tracking-[0.28em] text-muted animate-entrance">Quick actions</p>
            <h2 className="mt-2 text-xl font-semibold text-foreground sm:text-2xl animate-entrance-up">Move into the next banking tasks</h2>
          </div>
        </div>

        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-5 animate-stagger">
          {quickActions.map((action) => (
            <Card key={action.title} variant="interactive" className="h-full animate-hover-lift-strong">
              <CardContent className="flex h-full flex-col justify-between gap-4 p-5 sm:p-6">
                <div className="space-y-3">
                  <div className="flex items-center justify-between gap-3 animate-entrance">
                    <div className="flex h-11 w-11 items-center justify-center rounded-2xl border border-primary/30 bg-gradient-to-br from-primary/20 to-primary/5 text-sm font-semibold text-primary-light transition-all duration-300 group-hover:border-primary/50">
                      {action.title.slice(0, 1)}
                    </div>
                    <Badge variant="outline" isAnimated>{action.badge}</Badge>
                  </div>
                  <div className="animate-entrance-up">
                    <h3 className="text-base font-semibold text-foreground">{action.title}</h3>
                    <p className="mt-2 text-sm leading-6 text-muted">{action.description}</p>
                  </div>
                </div>
                <Button variant="secondary" size="sm" className="w-full justify-between animate-entrance-up" onClick={() => router.push(action.href)}>
                  <span>{action.badge}</span>
                  <span aria-hidden="true">↗</span>
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      </section>

      <Card variant="gradient" className="border-primary/20">
        <CardHeader className="gap-3 sm:flex-row sm:items-end sm:justify-between border-b border-border/40 bg-gradient-to-r from-primary/10 to-transparent">
          <div className="animate-entrance">
            <p className="text-sm uppercase tracking-[0.28em] text-muted animate-entrance">Recent activity</p>
            <CardTitle className="mt-2 text-xl text-white sm:text-2xl animate-entrance-up">Preview of transaction history</CardTitle>
            <CardDescription className="animate-entrance-up">
              Your recent account activity and transfers appear below in real-time.
            </CardDescription>
          </div>
          <Badge variant={activity.length > 0 ? 'status-active' : 'secondary'} isAnimated>{activity.length > 0 ? 'Live' : 'No activity'}</Badge>
        </CardHeader>
        <CardContent className="pt-6">
          {activity.length > 0 ? (
            <div className="overflow-hidden rounded-[1.25rem] border border-border/50 bg-gradient-to-b from-white/[0.04] to-transparent animate-entrance">
              <Table>
                <TableHeader className="bg-gradient-to-r from-primary/10 to-transparent border-b border-border/40">
                  <TableRow className="hover:bg-transparent">
                    <TableHead className="text-primary-light">Activity</TableHead>
                    <TableHead className="text-primary-light">Status</TableHead>
                    <TableHead className="text-primary-light">Details</TableHead>
                    <TableHead className="text-right text-primary-light">Time</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody className="animate-stagger">
                  {activity.map((event, index) => (
                    <TableRow key={`${event.type}-${index}`} className="border-border/30 hover:bg-primary/5 transition-all duration-200 animate-entrance-up">
                      <TableCell className="font-medium text-foreground">{getActivityTitle(event)}</TableCell>
                      <TableCell>
                        <Badge variant={getStatusBadgeVariant(event.status)} isAnimated>
                          {event.status.charAt(0).toUpperCase() + event.status.slice(1)}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-muted">{event.details}</TableCell>
                      <TableCell className="text-right text-muted">{formatTime(event.timestamp)}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          ) : (
            <div className="mt-5 animate-entrance-up">
              <EmptyState
                title="Transaction history will populate here"
                description="Transfers, offers, and passkey activity can all flow into this area without redesigning the shell."
                icon={<span className="text-2xl animate-pulse-smooth">⟡</span>}
              />
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
