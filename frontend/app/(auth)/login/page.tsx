'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { api, LoginRequest } from '@/lib/api';
import { auth } from '@/lib/auth';
import { Alert } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { AuthPage } from '@/components/auth-page';
import { PasskeyAlert } from '@/components/passkey-alert';
import {
  normalizeLoginOptions,
  normalizeLoginOptionsWithoutAllowList,
  serializeRequestCredential,
  unsupportedWebauthnMessage,
  webauthnSupported,
} from '@/lib/webauthn';

export default function LoginPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [passkeyLoading, setPasskeyLoading] = useState(false);
  const [error, setError] = useState('');
  const [passkeyError, setPasskeyError] = useState('');
  const [passkeyEmail, setPasskeyEmail] = useState('');
  const [formData, setFormData] = useState<LoginRequest>({
    number: '',
    password: '',
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const response = await api.login(formData);

      if (response.error) {
        setError(response.error);
        return;
      }

      if (response.data) {
        auth.setToken(response.data.token);
        auth.setUser({ number: response.data.number });
        router.push('/dashboard');
      }
    } catch {
      setError('An unexpected error occurred');
    } finally {
      setLoading(false);
    }
  };

  const handlePasskeyLogin = async () => {
    setPasskeyError('');

    if (!webauthnSupported()) {
      setPasskeyError(unsupportedWebauthnMessage());
      return;
    }

    if (!passkeyEmail.trim()) {
      setPasskeyError('Enter your email address so GoBank can start the passkey login ceremony.');
      return;
    }

    setPasskeyLoading(true);

    try {
      const beginResponse = await api.webauthnLoginBegin({ email: passkeyEmail.trim() });

      if (beginResponse.error || !beginResponse.data) {
        setPasskeyError(beginResponse.error || 'Unable to start passkey login');
        return;
      }

      let credential: Credential | null = null;

      try {
        credential = await navigator.credentials.get({
          publicKey: normalizeLoginOptions(beginResponse.data),
        });
      } catch (firstError) {
        const message = firstError instanceof Error ? firstError.message : '';
        const shouldRetry = message.toLowerCase().includes('allowed credential list') || message.toLowerCase().includes('credentials');

        if (!shouldRetry) {
          throw firstError;
        }

        credential = await navigator.credentials.get({
          publicKey: normalizeLoginOptionsWithoutAllowList(beginResponse.data),
        });
      }

      if (!credential || credential.type !== 'public-key') {
        setPasskeyError('Passkey login was cancelled before completion.');
        return;
      }

      const finishResponse = await api.webauthnLoginFinish(
        passkeyEmail.trim(),
        serializeRequestCredential(credential as PublicKeyCredential),
      );

      if (finishResponse.error || !finishResponse.data) {
        setPasskeyError(finishResponse.error || 'Passkey login failed');
        return;
      }

      auth.setToken(finishResponse.data.token);
      auth.setUser({ number: finishResponse.data.number });
      router.push('/dashboard');
    } catch (loginError) {
      setPasskeyError(loginError instanceof Error ? loginError.message : 'Passkey login failed');
    } finally {
      setPasskeyLoading(false);
    }
  };

  return (
    <AuthPage
      title="Welcome back"
      description="Sign in to your GoBank account and continue to your protected banking dashboard."
      eyebrow={
        <div className="flex flex-wrap gap-2 animate-entrance">
          <span className="rounded-full border border-primary/30 bg-primary/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.22em] text-primary-light animate-entrance-up">
            Secure login
          </span>
          <span className="rounded-full border border-accent/30 bg-accent/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.22em] text-accent-strong animate-entrance-up">
            JWT session
          </span>
        </div>
      }
      footer={
        <div className="space-y-4">
          <PasskeyAlert
            supported={webauthnSupported()}
            message={webauthnSupported()
              ? 'Passkey sign-in is available for enrolled users. Enter your account number and use a device passkey.'
              : unsupportedWebauthnMessage()}
          />
          {passkeyError ? <Alert variant="warning" title="Passkey login issue">{passkeyError}</Alert> : null}
          <p className="text-sm text-slate-300/80">
            Don&apos;t have an account?{' '}
            <Link href="/register" className="font-semibold text-accent hover:text-accent-strong">
              Sign up
            </Link>
          </p>
        </div>
      }
    >
      <form onSubmit={handleSubmit} className="space-y-5 animate-stagger">
        {error ? <Alert variant="destructive" title="Login error" icon="⚠️">{error}</Alert> : null}

        <div className="space-y-2 animate-entrance-up">
          <label htmlFor="number" className="text-sm font-medium text-slate-200">
            Account number
          </label>
          <Input
            id="number"
            name="number"
            type="text"
            inputMode="numeric"
            value={formData.number}
            onChange={handleChange}
            required
            placeholder="123456"
          />
        </div>

        <div className="space-y-2 animate-entrance-up">
          <label htmlFor="password" className="text-sm font-medium text-slate-200">
            Password
          </label>
          <Input
            id="password"
            name="password"
            type="password"
            value={formData.password}
            onChange={handleChange}
            required
            placeholder="Enter your password"
          />
        </div>

        <div className="space-y-2 rounded-2xl border border-primary/30 bg-gradient-to-r from-primary/10 to-transparent p-4 animate-entrance-up">
          <label htmlFor="passkeyEmail" className="text-sm font-medium text-slate-200">
            Passkey email
          </label>
          <Input
            id="passkeyEmail"
            name="passkeyEmail"
            type="email"
            value={passkeyEmail}
            onChange={(event) => setPasskeyEmail(event.target.value)}
            placeholder="you@example.com"
          />
          <p className="text-xs leading-6 text-slate-300/70">
            Use the email address linked to your enrolled passkey.
          </p>
        </div>

        <Button type="submit" disabled={loading} className="w-full animate-entrance-up">
          {loading ? 'Signing in...' : 'Sign In'}
        </Button>

        <div className="space-y-3 animate-entrance-up">
          <div className="relative flex items-center py-1">
            <div className="h-px flex-1 bg-gradient-to-r from-transparent via-primary/30 to-transparent" />
            <span className="px-3 text-[11px] font-semibold uppercase tracking-[0.28em] text-slate-400">or use a passkey</span>
            <div className="h-px flex-1 bg-gradient-to-r from-transparent via-primary/30 to-transparent" />
          </div>

          <Button type="button" variant="secondary" className="w-full hover:bg-surface-strong/80" onClick={handlePasskeyLogin} disabled={passkeyLoading}>
            {passkeyLoading ? 'Starting passkey login...' : 'Sign In with Passkey'}
          </Button>
        </div>
      </form>
    </AuthPage>
  );
}
