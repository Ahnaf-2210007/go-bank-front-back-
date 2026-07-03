'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { auth } from '@/lib/auth';
import { api, AccountResponse, UpdateProfileRequest } from '@/lib/api';
import { Alert } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { WorkflowPage } from '@/components/workflow-page';
import { PasskeyAlert } from '@/components/passkey-alert';
import {
  normalizeRegistrationOptions,
  serializeCreationCredential,
  unsupportedWebauthnMessage,
  webauthnSupported,
} from '@/lib/webauthn';

export default function ProfilePage() {
  const router = useRouter();
  const [account, setAccount] = useState<AccountResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [passkeyMessage, setPasskeyMessage] = useState('');
  const [passkeyError, setPasskeyError] = useState('');
  const [passkeyLoading, setPasskeyLoading] = useState(false);
  const [formData, setFormData] = useState<UpdateProfileRequest>({
    action: 'profile',
    firstName: '',
    lastName: '',
    newEmail: '',
    password: '',
    otp: '',
    currentPassword: '',
    newPassword: '',
    confirmPassword: '',
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
        setFormData((previous) => ({
          ...previous,
          firstName: response.data.firstName,
          lastName: response.data.lastName,
          newEmail: response.data.email,
        }));
      }

      setLoading(false);
    };

    loadAccount();

    return () => {
      mounted = false;
    };
  }, [router]);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = event.target;
    setFormData((previous) => ({
      ...previous,
      [name]: value,
    }));
  };

  const validate = () => {
    if (formData.action === 'profile') {
      if (!formData.firstName?.trim() || !formData.lastName?.trim()) {
        return 'First name and last name are required';
      }
    }

    if (formData.action === 'email_request') {
      if (!formData.newEmail?.trim() || !formData.password?.trim()) {
        return 'New email and password are required';
      }
    }

    if (formData.action === 'email_verify') {
      if (!formData.newEmail?.trim() || !formData.otp?.trim()) {
        return 'New email and 6-digit OTP are required';
      }
    }

    if (formData.action === 'password') {
      if (!formData.currentPassword?.trim() || !formData.newPassword?.trim() || !formData.confirmPassword?.trim()) {
        return 'All password fields are required';
      }
      if (formData.newPassword !== formData.confirmPassword) {
        return 'New password and confirmation must match';
      }
    }

    return '';
  };

  const submit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setMessage('');
    setError('');

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

    setSaving(true);
    try {
      const response = await api.updateAccount(token, formData);
      if (response.error) {
        setError(response.error);
        return;
      }

      if (response.data) {
        if ('email' in response.data) {
          setAccount(response.data);
          setMessage('Profile updated successfully.');
        } else {
          setMessage(response.data.message);
        }

        const refreshed = await api.getAccount(token);
        if (refreshed.data) {
          setAccount(refreshed.data);
        }
      }
    } finally {
      setSaving(false);
    }
  };

  const handlePasskeyRegistration = async () => {
    setPasskeyMessage('');
    setPasskeyError('');

    if (!webauthnSupported()) {
      setPasskeyError(unsupportedWebauthnMessage());
      return;
    }

    const token = auth.getToken();
    if (!token || !account?.email) {
      router.replace('/login');
      return;
    }

    setPasskeyLoading(true);

    try {
      const beginResponse = await api.webauthnRegisterBegin(token, { email: account.email });

      if (beginResponse.error || !beginResponse.data) {
        setPasskeyError(beginResponse.error || 'Unable to start passkey registration');
        return;
      }

      const credential = await navigator.credentials.create({
        publicKey: normalizeRegistrationOptions(beginResponse.data),
      });

      if (!credential || credential.type !== 'public-key') {
        setPasskeyError('Passkey registration was cancelled before completion.');
        return;
      }

      const finishResponse = await api.webauthnRegisterFinish(
        token,
        serializeCreationCredential(credential as PublicKeyCredential),
      );

      if (finishResponse.error || !finishResponse.data) {
        setPasskeyError(finishResponse.error || 'Passkey registration failed');
        return;
      }

      setPasskeyMessage('Passkey registered successfully. You can now use it from the login page.');
    } catch (registrationError) {
      setPasskeyError(registrationError instanceof Error ? registrationError.message : 'Passkey registration failed');
    } finally {
      setPasskeyLoading(false);
    }
  };

  if (loading) {
    return <Card><CardContent className="p-6 text-sm text-muted">Loading profile settings...</CardContent></Card>;
  }

  return (
    <WorkflowPage
      eyebrow="Profile"
      title="Update account details"
      description="Use the action-based backend contract to update your profile, request an email change, verify email, or update your password."
      primaryAction={
        <Button type="submit" form="profile-form" disabled={saving}>
          {saving ? 'Saving...' : 'Save changes'}
        </Button>
      }
      secondaryAction={
        <Button type="button" variant="secondary" onClick={() => router.push('/dashboard')}>
          Back to dashboard
        </Button>
      }
      summary={
        <div className="grid gap-4 md:grid-cols-3">
          <Card className="border-white/8 bg-white/[0.03]">
            <CardHeader className="pb-3">
              <CardDescription>Account holder</CardDescription>
              <CardTitle className="text-xl text-white">{account ? `${account.firstName} ${account.lastName}` : ''}</CardTitle>
            </CardHeader>
          </Card>
          <Card className="border-white/8 bg-white/[0.03]">
            <CardHeader className="pb-3">
              <CardDescription>Account number</CardDescription>
              <CardTitle className="text-xl tracking-[0.14em] text-white">{account?.number}</CardTitle>
            </CardHeader>
          </Card>
          <Card className="border-accent/15 bg-[linear-gradient(180deg,rgba(78,162,255,0.14),rgba(255,255,255,0.03))]">
            <CardHeader className="pb-3">
              <CardDescription>Email</CardDescription>
              <CardTitle className="text-lg text-white">{account?.email}</CardTitle>
            </CardHeader>
          </Card>
        </div>
      }
    >
      <div className="space-y-6">
        <Card className="border-accent/15 bg-[linear-gradient(180deg,rgba(78,162,255,0.12),rgba(255,255,255,0.03))]">
          <CardHeader className="pb-3">
            <CardDescription>Passkey enrollment</CardDescription>
            <CardTitle className="text-xl text-white">Register a device passkey</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4 pt-0">
            <PasskeyAlert
              supported={webauthnSupported()}
              message={webauthnSupported()
                ? 'Use a passkey to speed up future sign-ins without changing your password workflow.'
                : unsupportedWebauthnMessage()}
            />
            {passkeyError ? <Alert variant="warning" title="Passkey registration issue">{passkeyError}</Alert> : null}
            {passkeyMessage ? <Alert variant="success" title="Passkey registered">{passkeyMessage}</Alert> : null}
            <Button type="button" variant="secondary" onClick={handlePasskeyRegistration} disabled={passkeyLoading}>
              {passkeyLoading ? 'Registering passkey...' : 'Register Passkey'}
            </Button>
          </CardContent>
        </Card>

        <div className="mb-6 grid gap-3 md:grid-cols-3">
          <div className="rounded-2xl border border-border/60 bg-white/[0.03] p-4">
            <p className="text-xs uppercase tracking-[0.22em] text-muted">Account holder</p>
            <p className="mt-2 text-base font-semibold text-white">{account ? `${account.firstName} ${account.lastName}` : ''}</p>
          </div>
          <div className="rounded-2xl border border-border/60 bg-white/[0.03] p-4">
            <p className="text-xs uppercase tracking-[0.22em] text-muted">Account number</p>
            <p className="mt-2 text-base font-semibold tracking-[0.14em] text-white">{account?.number}</p>
          </div>
          <div className="rounded-2xl border border-border/60 bg-white/[0.03] p-4">
            <p className="text-xs uppercase tracking-[0.22em] text-muted">Email</p>
            <p className="mt-2 text-base font-semibold text-white">{account?.email}</p>
          </div>
        </div>

        <form id="profile-form" onSubmit={submit} className="space-y-5">
          {error ? <Alert variant="destructive" title="Update error">{error}</Alert> : null}
          {message ? <Alert variant="success" title="Update status">{message}</Alert> : null}

          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <label className="text-sm font-medium text-slate-200">Action</label>
              <select
                name="action"
                value={formData.action}
                onChange={handleChange}
                className="h-11 w-full rounded-2xl border border-border/70 bg-surface-strong/85 px-4 text-sm text-foreground outline-none transition focus:border-accent/70 focus:ring-2 focus:ring-accent/20"
              >
                <option value="profile">Profile name update</option>
                <option value="email_request">Email change request</option>
                <option value="email_verify">Email change verification</option>
                <option value="password">Password update</option>
              </select>
            </div>

            {formData.action === 'profile' ? (
              <div className="grid gap-4 md:grid-cols-2 md:col-span-1">
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-200">First name</label>
                  <Input name="firstName" value={formData.firstName || ''} onChange={handleChange} />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-200">Last name</label>
                  <Input name="lastName" value={formData.lastName || ''} onChange={handleChange} />
                </div>
              </div>
            ) : null}

            {formData.action === 'email_request' ? (
              <div className="grid gap-4 md:grid-cols-2 md:col-span-1">
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-200">New email</label>
                  <Input name="newEmail" type="email" value={formData.newEmail || ''} onChange={handleChange} />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-200">Current password</label>
                  <Input name="password" type="password" value={formData.password || ''} onChange={handleChange} />
                </div>
              </div>
            ) : null}

            {formData.action === 'email_verify' ? (
              <div className="grid gap-4 md:grid-cols-2 md:col-span-1">
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-200">New email</label>
                  <Input name="newEmail" type="email" value={formData.newEmail || ''} onChange={handleChange} />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-200">OTP</label>
                  <Input name="otp" inputMode="numeric" value={formData.otp || ''} onChange={handleChange} placeholder="123456" />
                </div>
              </div>
            ) : null}

            {formData.action === 'password' ? (
              <div className="grid gap-4 md:grid-cols-3 md:col-span-1">
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-200">Current password</label>
                  <Input name="currentPassword" type="password" value={formData.currentPassword || ''} onChange={handleChange} />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-200">New password</label>
                  <Input name="newPassword" type="password" value={formData.newPassword || ''} onChange={handleChange} />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-200">Confirm password</label>
                  <Input name="confirmPassword" type="password" value={formData.confirmPassword || ''} onChange={handleChange} />
                </div>
              </div>
            ) : null}

          <div className="flex flex-wrap items-center justify-end gap-3 border-t border-white/8 pt-4">
            <Button type="button" variant="ghost" onClick={() => router.push('/dashboard')} disabled={saving}>
              Cancel
            </Button>
            <Button type="submit" disabled={saving}>
              {saving ? 'Saving...' : 'Save changes'}
            </Button>
          </div>
          </div>

        </form>
      </div>
    </WorkflowPage>
  );
}
