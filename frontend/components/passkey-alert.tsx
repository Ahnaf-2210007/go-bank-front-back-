import { Alert } from '@/components/ui/alert';

interface PasskeyAlertProps {
  message: string;
  supported: boolean;
}

export function PasskeyAlert({ message, supported }: PasskeyAlertProps) {
  return (
    <Alert variant={supported ? 'default' : 'warning'} title={supported ? 'Passkey ready' : 'Passkey unavailable'}>
      {message}
    </Alert>
  );
}
