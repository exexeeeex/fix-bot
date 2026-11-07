import type { Application } from '@/entities/application';
import { Button } from '@/shared/components/ui/button';
import { useAcceptApplication, useApplicationCard } from '../../hooks';

interface ApplicationCardProps {
  application: Application;
}

export const ApplicationCard = ({ application }: ApplicationCardProps) => {
  const { statusLabel, statusClassName, typeLabel, createdAtLabel } =
    useApplicationCard(application);
  const { accept, isLoading, isError } = useAcceptApplication();

  return (
    <article className="rounded-xl border border-border bg-card p-4 text-card-foreground shadow-sm">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div className="min-w-0">
          <h3 className="break-words text-base font-semibold">{application.name}</h3>
          <p className="mt-1 text-sm text-muted-foreground">{typeLabel}</p>
        </div>
        <span
          className={`shrink-0 rounded-full px-2.5 py-1 text-xs font-medium ${statusClassName}`}
        >
          {statusLabel}
        </span>
      </div>

      <dl className="mt-4 grid gap-2 border-t border-border pt-3 text-sm sm:grid-cols-2">
        <div>
          <dt className="text-muted-foreground">Группа</dt>
          <dd className="mt-0.5 break-words font-medium">
            {application.group?.name ?? 'Не указана'}
          </dd>
        </div>
        <div>
          <dt className="text-muted-foreground">ID пользователя</dt>
          <dd className="mt-0.5 break-all font-medium">{application.userId}</dd>
        </div>
        {createdAtLabel && (
          <div>
            <dt className="text-muted-foreground">Подана</dt>
            <dd className="mt-0.5 font-medium">{createdAtLabel}</dd>
          </div>
        )}
      </dl>
      {application.status === 'pending' && (
        <div className="mt-4 border-t border-border pt-3">
          <Button
            type="button"
            className="w-full sm:w-auto"
            disabled={isLoading}
            onClick={() => void accept(application.id)}
          >
            {isLoading ? 'Принимаем...' : 'Принять заявку'}
          </Button>
          {isError && (
            <p role="alert" className="mt-2 text-sm text-destructive">
              Не удалось принять заявку. Попробуйте ещё раз.
            </p>
          )}
        </div>
      )}
    </article>
  );
};
