'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { auth } from '@/lib/auth';
import { api, AccountResponse, TransferRequest } from '@/lib/api';
import { Alert } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { WorkflowPage } from '@/components/workflow-page';

export default function TransferPage() {
  const router = useRouter();
  const [account, setAccount] = useState<AccountResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [formData, setFormData] = useState<TransferRequest>({
    toAccount: 0,
    amount: 0,
  });

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

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = event.target;
    setFormData((previous) => ({
      ...previous,
      [name]: Number(value),
    }));
  };

  const validate = () => {
    if (!formData.toAccount || Number.isNaN(formData.toAccount)) {
      return 'Recipient account number is required';
    }
    if (!formData.amount || Number.isNaN(formData.amount) || formData.amount <= 0) {
      return 'Enter a transfer amount greater than 0';
    }
    if (account && formData.toAccount === account.number) {
      return 'You cannot transfer to the same account';
    }
    return '';
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError('');
    setSuccess('');

    const validationError = validate();
    if (validationError) {
      setError(validationError);
      return;
    }

    const token = auth.getToken();
    if (!token) {
      router.replace('/login');
      return;
    }

    setSubmitting(true);

    try {
      const response = await api.transfer(token, formData);

      if (response.error) {
        setError(response.error);
        return;
      }

      if (response.data) {
        setSuccess(`Transfer completed successfully. Transaction ${response.data.transactionId}.`);
        setFormData({ toAccount: 0, amount: 0 });
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
    return (
      <WorkflowPage eyebrow="Payments" title="Transfer funds" description="Move money securely between GoBank accounts." summary={null}>
        <div className="text-sm text-muted">Loading transfer form...</div>
      </WorkflowPage>
    );
  }

  return (
    <WorkflowPage
      eyebrow="Payments"
      title="Transfer funds"
      description="Move money securely between GoBank accounts with clear validation and instant feedback."
      primaryAction={
        <Button type="submit" form="transfer-form" disabled={submitting}>
          {submitting ? 'Sending...' : 'Send transfer'}
        </Button>
      }
      secondaryAction={
        <Button type="button" variant="secondary" onClick={() => router.push('/dashboard')}>
          Back to dashboard
        </Button>
      }
      summary={
        <div className="grid gap-4 md:grid-cols-3">
          <Card className="border-accent/15 bg-white/[0.03]">
            <CardHeader className="pb-3">
              <CardDescription>Current balance</CardDescription>
              <CardTitle className="text-3xl text-white">${account?.balance.toFixed(2)}</CardTitle>
            </CardHeader>
          </Card>
          <Card className="border-white/8 bg-white/[0.03]">
            <CardHeader className="pb-3">
              <CardDescription>Account number</CardDescription>
              <CardTitle className="text-xl tracking-[0.14em] text-white">{account?.number}</CardTitle>
            </CardHeader>
          </Card>
          <Card className="border-white/8 bg-white/[0.03]">
            <CardHeader className="pb-3">
              <CardDescription>Account holder</CardDescription>
              <CardTitle className="text-xl text-white">{account ? `${account.firstName} ${account.lastName}` : ''}</CardTitle>
            </CardHeader>
          </Card>
        </div>
      }
    >
      <div className="space-y-5">
        <div className="flex flex-wrap items-center gap-2">
          <Badge variant="default">Transfer funds</Badge>
          <Badge variant="outline">Positive amounts only</Badge>
        </div>

        <div className="space-y-2">
          <h2 className="text-2xl font-semibold tracking-tight text-white">Send money securely</h2>
          <p className="max-w-2xl text-sm leading-6 text-slate-200/75">
            Transfer to another GoBank account using the backend transfer endpoint.
          </p>
        </div>

        <form id="transfer-form" onSubmit={handleSubmit} className="space-y-5">
          {error ? <Alert variant="destructive" title="Transfer error">{error}</Alert> : null}
          {success ? <Alert variant="success" title="Transfer complete">{success}</Alert> : null}

          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <label htmlFor="toAccount" className="text-sm font-medium text-slate-200">Recipient account number</label>
              <Input
                id="toAccount"
                name="toAccount"
                type="number"
                inputMode="numeric"
                value={formData.toAccount || ''}
                onChange={handleChange}
                placeholder="987654"
              />
            </div>

            <div className="space-y-2">
              <label htmlFor="amount" className="text-sm font-medium text-slate-200">Amount</label>
              <Input
                id="amount"
                name="amount"
                type="number"
                inputMode="decimal"
                value={formData.amount || ''}
                onChange={handleChange}
                placeholder="1000"
              />
            </div>
          </div>
        </form>
      </div>
    </WorkflowPage>
  );
}
