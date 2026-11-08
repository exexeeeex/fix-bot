import { Icon } from '@/shared/components';
import { ApplicationStatusCard } from './application-status-card';

export const PendingApplicationCard = () => (
  <ApplicationStatusCard
    title="Заявка на рассмотрении"
    description="Мы получили вашу заявку. Она будет доступна после проверки администрацией."
    tone="pending"
    icon={<Icon icon="clock" aria-hidden className="[&_svg]:size-8" />}
  />
);
