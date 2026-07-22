'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { auth } from '@/lib/auth';
import { api, AccountResponse } from '@/lib/api';
import { Alert } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { WorkflowPage } from '@/components/workflow-page';

export default function CouponPage() {
  const router = useRouter();
  const [account, setAccount] = useState<AccountResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [couponCode, setCouponCode] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

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

  const submit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError('');
    setSuccess('');

    if (!couponCode.trim()) {
      setError('Coupon code is required');
      return;
    }

    const token = auth.getToken();
    if (!token || !account) {
      router.replace('/login');
      return;
    }

    setSubmitting(true);
    try {
      const response = await api.redeemCoupon(token, account.id, couponCode.trim());
      if (response.error) {
        setError(response.error);
        return;
      }

      if (response.data) {
        setSuccess(response.data.status);
        setCouponCode('');
        const refreshed = await api.getAccount(token);
        if (refreshed.data) {
          setAccount(refreshed.data);
        }
      }
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return <Card><CardContent className="p-6 text-sm text-muted">Loading coupon center...</CardContent></Card>;
  }

  return (
    <WorkflowPage
      eyebrow="Offers"
      title="Coupon redemption"
      description="Redeem a coupon or offer using the backend redemption endpoint and refresh your balance after success."
      primaryAction={
        <Button type="submit" form="coupon-form" disabled={submitting}>
          {submitting ? 'Applying...' : 'Apply coupon'}
        </Button>
      }
      secondaryAction={
        <Button type="button" variant="secondary" onClick={() => router.push('/dashboard')}>
          Back to dashboard
        </Button>
      }
      summary={
        <div className="grid gap-4 md:grid-cols-3 animate-stagger">
          <Card variant="default" className="animate-hover-lift">
            <CardHeader className="pb-3 border-b border-border/40 bg-gradient-to-r from-accent-secondary/10 to-transparent">
              <CardDescription>Account holder</CardDescription>
              <CardTitle className="text-xl text-accent-strong animate-entrance-up">{account ? `${account.firstName} ${account.lastName}` : ''}</CardTitle>
            </CardHeader>
          </Card>
          <Card variant="elevated" className="animate-hover-lift">
            <CardHeader className="pb-3 border-b border-border/40 bg-gradient-to-r from-primary/10 to-transparent">
              <CardDescription>Account number</CardDescription>
              <CardTitle className="text-xl tracking-[0.14em] text-primary-light animate-entrance-up">{account?.number}</CardTitle>
            </CardHeader>
          </Card>
          <Card variant="featured" className="animate-hover-lift-strong">
            <CardHeader className="pb-3 border-b border-border/40 bg-gradient-to-r from-success/10 to-transparent">
              <CardDescription>Current balance</CardDescription>
              <CardTitle className="text-3xl text-success-light animate-entrance-up">${account?.balance.toFixed(2)}</CardTitle>
            </CardHeader>
          </Card>
        </div>
      }
    >
      <div className="space-y-5 animate-stagger">
        <div className="grid gap-3 md:grid-cols-3 animate-entrance">
          <div className="rounded-2xl border border-border/60 bg-gradient-to-br from-white/[0.04] to-transparent p-4 hover:bg-white/[0.06] transition-colors duration-300 animate-entrance-up">
            <p className="text-xs uppercase tracking-[0.22em] text-muted">Account holder</p>
            <p className="mt-2 text-base font-semibold text-white">{account ? `${account.firstName} ${account.lastName}` : ''}</p>
          </div>
          <div className="rounded-2xl border border-border/60 bg-gradient-to-br from-white/[0.04] to-transparent p-4 hover:bg-white/[0.06] transition-colors duration-300 animate-entrance-up">
            <p className="text-xs uppercase tracking-[0.22em] text-muted">Account number</p>
            <p className="mt-2 text-base font-semibold tracking-[0.14em] text-white">{account?.number}</p>
          </div>
          <div className="rounded-2xl border border-success/20 bg-gradient-to-br from-success/12 to-success/5 p-4 hover:from-success/14 hover:to-success/8 transition-all duration-300 animate-entrance-up">
            <p className="text-xs uppercase tracking-[0.22em] text-slate-200/70">Current balance</p>
            <p className="mt-2 text-2xl font-bold text-success-light">${account?.balance.toFixed(2)}</p>
          </div>
        </div>

        {error ? <Alert variant="destructive" title="Coupon error" icon="✕">{error}</Alert> : null}
        {success ? <Alert variant="success" title="Coupon applied" icon="✓">{success}</Alert> : null}

        <form id="coupon-form" onSubmit={submit} className="space-y-4 animate-entrance-up">
          <div className="space-y-2">
            <label htmlFor="couponCode" className="text-sm font-medium text-slate-200">Coupon code</label>
            <Input
              id="couponCode"
              value={couponCode}
              onChange={(event) => setCouponCode(event.target.value)}
              placeholder="BANK1000"
            />
          </div>

        </form>
      </div>
    </WorkflowPage>
  );
}
