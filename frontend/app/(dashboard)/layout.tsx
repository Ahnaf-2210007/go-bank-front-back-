'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { auth } from '@/lib/auth';
import { DashboardShell } from '@/components/dashboard-shell';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const [ready, setReady] = useState(false);

  useEffect(() => {
    if (!auth.isAuthenticated()) {
      router.replace('/login');
      return;
    }

    setReady(true);
  }, [router]);

  if (!ready) {
    return (
      <div className="gobank-shell items-center justify-center px-4 py-8">
        <Card className="w-full max-w-lg">
          <CardHeader className="space-y-4">
            <Skeleton className="h-4 w-28" />
            <Skeleton className="h-9 w-60 max-w-full" />
            <Skeleton className="h-5 w-80 max-w-full" />
          </CardHeader>
          <CardContent className="space-y-3">
            <Skeleton className="h-12 w-full rounded-2xl" />
            <Skeleton className="h-12 w-3/4 rounded-2xl" />
          </CardContent>
        </Card>
      </div>
    );
  }

  return <DashboardShell>{children}</DashboardShell>;
}
