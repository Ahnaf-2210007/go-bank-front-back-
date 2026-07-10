'use client';

import { useState, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import Link from 'next/link';
import { api } from '@/lib/api';
import { Alert } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { AuthPage } from '@/components/auth-page';

function VerifyContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [code, setCode] = useState('');
  const email = searchParams.get('email') || '';

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setLoading(true);

    try {
      const response = await api.verifyEmail({ code });

      if (response.error) {
        setError(response.error);
        return;
      }

      if (response.data) {
        setSuccess('Email verified successfully! Redirecting to login...');
        setTimeout(() => {
          router.push('/login');
        }, 2000);
      }
    } catch {
      setError('An unexpected error occurred');
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthPage
      title="Verify email"
      description={`Enter the verification code sent to ${email || 'your email'}.`}
      eyebrow={
        <div className="flex flex-wrap gap-2 animate-entrance">
          <span className="rounded-full border border-success/30 bg-success/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.22em] text-success-light animate-entrance-up">
            Email verification
          </span>
          <span className="rounded-full border border-accent/30 bg-accent/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.22em] text-accent-strong animate-entrance-up">
            6 digits
          </span>
        </div>
      }
      footer={
        <p className="text-sm text-slate-300/80">
          Already verified?{' '}
          <Link href="/login" className="font-semibold text-accent hover:text-accent-strong">
            Go to login
          </Link>
        </p>
      }
    >
      <form onSubmit={handleSubmit} className="space-y-5 animate-stagger">
        {error ? <Alert variant="destructive" title="Verification error" icon="✕">{error}</Alert> : null}
        {success ? <Alert variant="success" title="Verification complete" icon="✓">{success}</Alert> : null}

        <div className="space-y-2 animate-entrance-up">
          <label htmlFor="code" className="text-sm font-medium text-slate-200">Verification code</label>
          <Input
            id="code"
            type="text"
            value={code}
            onChange={(e) => setCode(e.target.value)}
            required
            maxLength={6}
            className="text-center text-2xl tracking-[0.35em] bg-gradient-to-r from-primary/10 to-transparent border-primary/40"
            placeholder="000000"
          />
        </div>

        <p className="text-xs leading-6 text-slate-300/70 animate-entrance-up">
          Check your email for the verification code. It may take a few minutes to arrive.
        </p>

        <Button type="submit" disabled={loading} className="w-full animate-entrance-up">
          {loading ? 'Verifying...' : 'Verify Email'}
        </Button>
      </form>
    </AuthPage>
  );
}

export default function VerifyPage() {
  return (
    <Suspense fallback={
      <div className="flex min-h-screen items-center justify-center bg-background px-4 text-muted">
        Loading...
      </div>
    }>
      <VerifyContent />
    </Suspense>
  );
}
