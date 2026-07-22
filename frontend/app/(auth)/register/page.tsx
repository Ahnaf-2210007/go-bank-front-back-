'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { api, RegisterRequest } from '@/lib/api';
import { Alert } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { AuthPage } from '@/components/auth-page';

export default function RegisterPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [formData, setFormData] = useState<RegisterRequest>({
    email: '',
    password: '',
    firstName: '',
    lastName: '',
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
    setSuccess('');
    setLoading(true);

    try {
      const response = await api.register(formData);

      if (response.error) {
        setError(response.error);
        return;
      }

      if (response.data) {
        setSuccess('Account created successfully! Redirecting to verification...');
        setTimeout(() => {
          router.push(`/verify?email=${encodeURIComponent(formData.email)}`);
        }, 1500);
      }
    } catch {
      setError('An unexpected error occurred');
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthPage
      title="Create account"
      description="Join GoBank to set up your secure banking profile and continue to email verification."
      eyebrow={
        <div className="flex flex-wrap gap-2 animate-entrance">
          <span className="rounded-full border border-featured/30 bg-featured/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.22em] text-featured-light animate-entrance-up">
            New customer
          </span>
          <span className="rounded-full border border-success/30 bg-success/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.22em] text-success-light animate-entrance-up">
            Secure signup
          </span>
        </div>
      }
      footer={
        <p className="text-sm text-slate-300/80">
          Already have an account?{' '}
          <Link href="/login" className="font-semibold text-accent hover:text-accent-strong">
            Sign in
          </Link>
        </p>
      }
    >
      <form onSubmit={handleSubmit} className="space-y-5 animate-stagger">
        {error ? <Alert variant="destructive" title="Registration error" icon="✕">{error}</Alert> : null}
        {success ? <Alert variant="success" title="Account created" icon="✓">{success}</Alert> : null}

        <div className="grid gap-4 sm:grid-cols-2 animate-entrance">
          <div className="space-y-2 animate-entrance-up">
            <label htmlFor="firstName" className="text-sm font-medium text-slate-200">First name</label>
            <Input id="firstName" name="firstName" type="text" value={formData.firstName} onChange={handleChange} required placeholder="John" />
          </div>

          <div className="space-y-2 animate-entrance-up">
            <label htmlFor="lastName" className="text-sm font-medium text-slate-200">Last name</label>
            <Input id="lastName" name="lastName" type="text" value={formData.lastName} onChange={handleChange} required placeholder="Doe" />
          </div>
        </div>

        <div className="space-y-2 animate-entrance-up">
          <label htmlFor="email" className="text-sm font-medium text-slate-200">Email</label>
          <Input id="email" name="email" type="email" value={formData.email} onChange={handleChange} required placeholder="you@example.com" />
        </div>

        <div className="space-y-2 animate-entrance-up">
          <label htmlFor="password" className="text-sm font-medium text-slate-200">Password</label>
          <Input id="password" name="password" type="password" value={formData.password} onChange={handleChange} required placeholder="Create a password" />
        </div>

        <Button type="submit" disabled={loading} className="w-full animate-entrance-up">
          {loading ? 'Creating Account...' : 'Create Account'}
        </Button>
      </form>
    </AuthPage>
  );
}
