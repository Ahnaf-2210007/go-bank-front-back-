'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { auth } from '@/lib/auth';
import { api, TransactionRecord } from '@/lib/api';
import { Alert } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { EmptyState } from '@/components/ui/empty-state';
import { Input } from '@/components/ui/input';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { WorkflowPage } from '@/components/workflow-page';

const monthOptions = [
  { label: 'Any month', value: '' },
  { label: 'January', value: 'january' },
  { label: 'February', value: 'february' },
  { label: 'March', value: 'march' },
  { label: 'April', value: 'april' },
  { label: 'May', value: 'may' },
  { label: 'June', value: 'june' },
  { label: 'July', value: 'july' },
  { label: 'August', value: 'august' },
  { label: 'September', value: 'september' },
  { label: 'October', value: 'october' },
  { label: 'November', value: 'november' },
  { label: 'December', value: 'december' },
];

export default function HistoryPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [records, setRecords] = useState<TransactionRecord[]>([]);
  const [error, setError] = useState('');
  const [filters, setFilters] = useState({
    type: '',
    month: '',
    limit: 10,
    offset: 0,
  });

  const loadHistory = async (nextFilters = filters) => {
    const token = auth.getToken();
    if (!token) {
      router.replace('/login');
      return;
    }

    setLoading(true);
    setError('');

    const response = await api.getTransactions(token, {
      limit: nextFilters.limit,
      offset: nextFilters.offset,
      type: nextFilters.type || undefined,
      month: nextFilters.month || undefined,
    });

    if (response.error) {
      setError(response.error);
      setRecords([]);
    } else {
      setRecords(response.data ?? []);
    }

    setLoading(false);
  };

  useEffect(() => {
    void loadHistory();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const setFilterValue = (name: 'type' | 'month' | 'limit', value: string) => {
    setFilters((previous) => ({
      ...previous,
      [name]: name === 'limit' ? Number(value) : value,
      offset: 0,
    }));
  };

  const applyFilters = async () => {
    await loadHistory({ ...filters, offset: 0 });
  };

  const goToPage = async (direction: 'next' | 'prev') => {
    const nextOffset = direction === 'next' ? filters.offset + filters.limit : Math.max(0, filters.offset - filters.limit);
    const nextFilters = { ...filters, offset: nextOffset };
    setFilters(nextFilters);
    await loadHistory(nextFilters);
  };

  return (
    <WorkflowPage
      eyebrow="Activity"
      title="Transaction history"
      description="Filter activity by type or month and page through readable transaction records."
      primaryAction={
        <Button onClick={applyFilters} disabled={loading}>
          Apply filters
        </Button>
      }
      secondaryAction={
        <Button variant="secondary" onClick={() => router.push('/dashboard')}>
          Back to dashboard
        </Button>
      }
      summary={
        <div className="grid gap-4 md:grid-cols-3">
          <Card className="border-white/8 bg-white/[0.03]">
            <CardHeader className="pb-3">
              <CardDescription>Rows per page</CardDescription>
              <CardTitle className="text-3xl text-white">{filters.limit}</CardTitle>
            </CardHeader>
          </Card>
          <Card className="border-white/8 bg-white/[0.03]">
            <CardHeader className="pb-3">
              <CardDescription>Current offset</CardDescription>
              <CardTitle className="text-3xl text-white">{filters.offset}</CardTitle>
            </CardHeader>
          </Card>
          <Card className="border-accent/15 bg-[linear-gradient(180deg,rgba(78,162,255,0.14),rgba(255,255,255,0.03))]">
            <CardHeader className="pb-3">
              <CardDescription>Loaded records</CardDescription>
              <CardTitle className="text-3xl text-white">{records.length}</CardTitle>
            </CardHeader>
          </Card>
        </div>
      }
    >
      <div className="space-y-5">
        {error ? <Alert variant="destructive" title="History error">{error}</Alert> : null}

        <div className="grid gap-4 lg:grid-cols-[1.4fr_1fr_1fr_1fr_auto]">
          <div className="space-y-2">
            <label htmlFor="type" className="text-sm font-medium text-slate-200">Transaction type</label>
            <Input
              id="type"
              value={filters.type}
              onChange={(event) => setFilterValue('type', event.target.value)}
              placeholder="transfer, coupon, deposit..."
            />
          </div>

          <div className="space-y-2">
            <label htmlFor="month" className="text-sm font-medium text-slate-200">Month</label>
            <select
              id="month"
              value={filters.month}
              onChange={(event) => setFilterValue('month', event.target.value)}
              className="h-11 w-full rounded-2xl border border-border/70 bg-surface-strong/85 px-4 text-sm text-foreground outline-none transition focus:border-accent/70 focus:ring-2 focus:ring-accent/20"
            >
              {monthOptions.map((option) => (
                <option key={option.label} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>

          <div className="space-y-2">
            <label htmlFor="limit" className="text-sm font-medium text-slate-200">Rows per page</label>
            <select
              id="limit"
              value={filters.limit}
              onChange={(event) => setFilterValue('limit', event.target.value)}
              className="h-11 w-full rounded-2xl border border-border/70 bg-surface-strong/85 px-4 text-sm text-foreground outline-none transition focus:border-accent/70 focus:ring-2 focus:ring-accent/20"
            >
              {[5, 10, 20, 50, 100].map((option) => (
                <option key={option} value={option}>
                  {option}
                </option>
              ))}
            </select>
          </div>

          <div className="flex items-end gap-3 lg:justify-end">
            <Button onClick={applyFilters} disabled={loading}>Apply</Button>
            <Button variant="secondary" onClick={() => router.push('/dashboard')}>Back</Button>
          </div>
        </div>

        {loading ? (
          <div className="rounded-[1.25rem] border border-border/60 p-6 text-sm text-muted">Loading transactions...</div>
        ) : records.length === 0 ? (
          <EmptyState
            title="No activity found"
            description="Try widening the filters or loading another month."
            icon={<span className="text-2xl">⌁</span>}
          />
        ) : (
          <div className="overflow-hidden rounded-[1.25rem] border border-border/60 bg-white/[0.02]">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Date</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Counterparty</TableHead>
                  <TableHead className="text-right">Amount</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {records.map((record) => {
                  const outgoing = record.transactionType === 'transfer' && record.fromAccountId !== record.toAccountId;
                  const counterparty = outgoing ? `To ${record.toAccountId}` : `From ${record.fromAccountId}`;

                  return (
                    <TableRow key={record.id}>
                      <TableCell className="whitespace-nowrap text-slate-200">{new Date(record.createdAt).toLocaleString()}</TableCell>
                      <TableCell className="capitalize text-foreground">{record.transactionType}</TableCell>
                      <TableCell>
                        <Badge variant={record.status === 'completed' || record.status === 'success' ? 'success' : 'secondary'}>
                          {record.status}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-muted">{counterparty}</TableCell>
                      <TableCell className={`text-right font-semibold ${outgoing ? 'text-danger' : 'text-success'}`}>
                        {outgoing ? '-' : '+'}${Number(record.amount).toFixed(2)}
                      </TableCell>
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>
          </div>
        )}

        <div className="flex flex-wrap items-center justify-between gap-3 border-t border-border/60 pt-4">
          <p className="text-sm text-muted">Offset {filters.offset} · limit {filters.limit}</p>
          <div className="flex gap-3">
            <Button variant="secondary" onClick={() => void goToPage('prev')} disabled={loading || filters.offset === 0}>
              Previous
            </Button>
            <Button onClick={() => void goToPage('next')} disabled={loading || records.length < filters.limit}>
              Next
            </Button>
          </div>
        </div>
      </div>
    </WorkflowPage>
  );
}
