import type { Route } from '../+types/root';
import { Home } from '@/widgets/home';

export function meta(_args: Route.MetaArgs) {
  return [
    { title: 'Фикс пропуск' },
    { name: 'description', content: 'Электронный пропуск студента' },
  ];
}

export default function HomePage() {
  return <Home />;
}
