import { Icon } from '@/shared/components';
import { Button } from '@/shared/components/ui/button';
import { ApplicationStatusCard } from './application-status-card';

type AuthenticationStatusCardProps = {
  status: 'loading' | 'error';
  onRetry: () => void;
};

export const AuthenticationStatusCard = ({ status, onRetry }: AuthenticationStatusCardProps) => {
  const isError = status === 'error';

  return (
    <ApplicationStatusCard
      title={isError ? 'Не удалось проверить статус' : 'Проверяем вашу заявку'}
      description={
        isError
          ? 'Проверьте подключение к интернету и попробуйте ещё раз.'
          : 'Подождите немного, загружаем информацию о вашей заявке.'
      }
      tone={isError ? 'warning' : 'info'}
      icon={
        isError ? (
          <Icon icon="alert-circle" aria-hidden className="[&_svg]:size-8" />
        ) : (
          <Icon icon="loader-circle" aria-hidden className="[&_svg]:size-8 [&_svg]:animate-spin" />
        )
      }
      action={
        isError ? (
          <Button
            type="button"
            onClick={onRetry}
            className="mt-6 h-12 rounded-xl px-5 text-base font-semibold"
          >
            <Icon icon="rotate-cw" aria-hidden className="mr-2 [&_svg]:size-4" />
            Проверить ещё раз
          </Button>
        ) : undefined
      }
    />
  );
};
