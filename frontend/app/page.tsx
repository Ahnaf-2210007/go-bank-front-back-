'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { auth } from '@/lib/auth';

export default function Home() {
  const router = useRouter();

  useEffect(() => {
    if (auth.isAuthenticated()) {
      router.push('/dashboard');
    } else {
      router.push('/login');
    }
  }, [router]);

  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-6">
      <div className="gobank-panel-strong w-full max-w-md space-y-4 p-8 text-center shadow-[0_26px_90px_-50px_rgba(78,162,255,0.7)]">
        <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl border border-accent/20 bg-accent/12 text-lg font-semibold text-accent">
          G
        </div>
        <div className="space-y-2">
          <p className="text-sm uppercase tracking-[0.28em] text-muted">GoBank</p>
          <h1 className="text-2xl font-semibold text-foreground">Restoring your secure session</h1>
          <p className="text-sm leading-6 text-slate-200/75">Checking your token and routing you to the right place.</p>
        </div>
        <div className="mx-auto h-12 w-12 animate-spin rounded-full border-2 border-white/12 border-t-accent" />
      </div>
    </div>
  );
}
