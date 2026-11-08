import { ApplicationsPage } from '@/features/application/ui/applications-page';
import type { Route } from '../+types/root';

export function meta(_args: Route.MetaArgs) {
  return [{ title: 'Заявки' }, { name: 'description', content: 'Список заявок пользователей' }];
}

export default function ApplicationsRoute() {
  return <ApplicationsPage />;
}
