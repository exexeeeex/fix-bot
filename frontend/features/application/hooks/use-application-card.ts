import type { Application } from '@/entities/application';

const statusLabels: Record<Application['status'], string> = {
  pending: 'На рассмотрении',
  approved: 'Одобрена',
  rejected: 'Отклонена',
  reject: 'Отклонена',
};

const statusStyles: Record<Application['status'], string> = {
  pending: 'bg-amber-100 text-amber-800',
  approved: 'bg-green-100 text-green-800',
  rejected: 'bg-red-100 text-red-800',
  reject: 'bg-red-100 text-red-800',
};

const applicationTypeLabels: Record<Application['type'], string> = {
  student: 'Студент',
  teacher: 'Преподаватель',
};

export const useApplicationCard = (application: Application) => {
  const timestamp = Date.parse(application.createdAt);

  return {
    statusLabel: statusLabels[application.status],
    statusClassName: statusStyles[application.status],
    typeLabel: applicationTypeLabels[application.type],
    createdAtLabel: Number.isFinite(timestamp)
      ? new Intl.DateTimeFormat('ru-RU', { dateStyle: 'medium' }).format(timestamp)
      : null,
  };
};
