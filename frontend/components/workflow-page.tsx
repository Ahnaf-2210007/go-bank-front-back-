import type { ReactNode } from 'react';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

interface WorkflowPageProps {
  eyebrow: string;
  title: string;
  description: string;
  primaryAction?: ReactNode;
  secondaryAction?: ReactNode;
  summary?: ReactNode;
  children: ReactNode;
}

export function WorkflowPage({
  eyebrow,
  title,
  description,
  primaryAction,
  secondaryAction,
  summary,
  children,
}: WorkflowPageProps) {
  return (
    <div className="space-y-6 lg:space-y-8">
      <Card className="overflow-hidden border-accent/12 bg-[linear-gradient(135deg,rgba(14,22,41,0.98),rgba(11,17,32,0.92))]">
        <CardContent className="relative p-6 sm:p-8 lg:p-10">
          <div className="absolute inset-0 -z-0 bg-[radial-gradient(circle_at_top_right,rgba(78,162,255,0.18),transparent_26%),radial-gradient(circle_at_bottom_left,rgba(57,217,138,0.1),transparent_24%)]" />
          <div className="relative z-10 flex flex-col gap-6 lg:flex-row lg:items-end lg:justify-between">
            <div className="space-y-4">
              <Badge variant="default">{eyebrow}</Badge>
              <div className="space-y-2">
                <CardTitle className="text-3xl tracking-tight text-white sm:text-4xl">{title}</CardTitle>
                <CardDescription className="max-w-2xl text-base leading-7 text-slate-200/80">{description}</CardDescription>
              </div>
            </div>

            {(primaryAction || secondaryAction) && (
              <div className="flex flex-wrap gap-3">
                {primaryAction}
                {secondaryAction}
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {summary}

      <Card className="border-white/8 bg-[linear-gradient(180deg,rgba(12,19,35,0.96),rgba(9,14,26,0.96))]">
        <CardContent className="p-6 sm:p-8">{children}</CardContent>
      </Card>
    </div>
  );
}
