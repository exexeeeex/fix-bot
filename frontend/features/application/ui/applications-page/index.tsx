import { useApplicationsPage } from '../../hooks';
import { ApplicationCard } from '../application-card';
import { Button } from '@/shared/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select';

const ALL_GROUPS = 'all';

const ApplicationsPage = () => {
  const {
    groupId,
    selectGroup,
    groupOptions,
    visibleGroups,
    pageState,
    isGroupsError,
    emptyMessage,
    retryApplications,
    retryGroups,
  } = useApplicationsPage();

  return (
    <main className="mx-auto min-h-screen w-full max-w-5xl px-4 py-8 pb-24 sm:px-6">
      <header className="mb-6">
        <h1 className="text-2xl font-semibold tracking-tight sm:text-3xl">Заявки</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          Просматривайте заявки вашей организации и фильтруйте их по группам.
        </p>
      </header>

      <section
        aria-label="Фильтры заявок"
        className="mb-6 rounded-xl border border-border bg-card p-4"
      >
        <div className="grid max-w-md gap-2">
          <label htmlFor="application-group" className="text-sm font-medium">
            Группа
          </label>
          <Select value={groupId} onValueChange={selectGroup} disabled={groupOptions.length === 0}>
            <SelectTrigger id="application-group" className="w-full">
              <SelectValue placeholder="Все группы" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value={ALL_GROUPS}>Все группы</SelectItem>
              {groupOptions.map((group) => (
                <SelectItem key={group.id} value={group.id}>
                  {group.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </section>

      {isGroupsError && pageState === 'content' && (
        <div className="mb-4 flex flex-wrap items-center justify-between gap-3 rounded-lg border border-amber-300 bg-amber-50 p-4 text-sm text-amber-900">
          <p>Не удалось загрузить список групп. Используются группы из заявок.</p>
          <Button type="button" variant="outline" onClick={() => void retryGroups()}>
            Повторить
          </Button>
        </div>
      )}

      {pageState === 'loading-organization' || pageState === 'loading-applications' ? (
        <p className="rounded-lg border border-border p-6 text-center text-sm text-muted-foreground">
          {pageState === 'loading-organization'
            ? 'Загружаем данные пользователя...'
            : 'Загружаем заявки...'}
        </p>
      ) : pageState === 'no-organization' ? (
        <EmptyMessage>
          Для вашей учётной записи не указана организация. Обратитесь к администратору.
        </EmptyMessage>
      ) : pageState === 'applications-error' ? (
        <ErrorMessage onRetry={() => void retryApplications()} />
      ) : visibleGroups.length === 0 ? (
        <EmptyMessage>{emptyMessage}</EmptyMessage>
      ) : (
        <div className="space-y-6">
          {visibleGroups.map((group) => (
            <section key={group.id} aria-labelledby={`group-${group.id}`}>
              <div className="mb-3 flex items-baseline justify-between gap-2">
                <h2 id={`group-${group.id}`} className="text-base font-semibold">
                  {group.name}
                </h2>
              </div>
              <div className="grid gap-3 md:grid-cols-2">
                {group.applications.map((application) => (
                  <ApplicationCard key={application.id} application={application} />
                ))}
              </div>
            </section>
          ))}
        </div>
      )}
    </main>
  );
};

const EmptyMessage = ({ children }: { children: React.ReactNode }) => (
  <p className="rounded-lg border border-dashed border-border p-8 text-center text-sm text-muted-foreground">
    {children}
  </p>
);

const ErrorMessage = ({ onRetry }: { onRetry: () => void }) => (
  <div
    role="alert"
    className="flex flex-col items-center gap-3 rounded-lg border border-destructive/30 bg-destructive/5 p-6 text-center"
  >
    <p className="text-sm text-destructive">Не удалось загрузить данные. Попробуйте ещё раз.</p>
    <Button type="button" variant="outline" onClick={onRetry}>
      Повторить
    </Button>
  </div>
);

export default ApplicationsPage;
