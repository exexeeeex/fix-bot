import { Icon } from '@/shared/components';
import { ApplicationStatusCard } from './application-status-card';

export const RejectedApplicationCard = () => (
  <ApplicationStatusCard
    title="Заявка отклонена"
    description="К сожалению, ваша заявка отклонена. Если вы считаете, что это произошло по ошибке, обратитесь к администратору."
    tone="danger"
    icon={<Icon icon="alert-circle" aria-hidden className="[&_svg]:size-8" />}
  />
);
