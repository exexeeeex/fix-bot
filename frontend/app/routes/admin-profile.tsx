import { Scanner } from '@/widgets/scanner';
import type { Route } from '../+types/root';

export function meta(_args: Route.MetaArgs) {
  return [{ title: 'Сканнер' }, { name: 'description', content: 'Сканнер для пропусков' }];
}

export default function AdminProfilePage() {
  return (
    <section className="min-h-screen">
      <Scanner />
    </section>
  );
}
