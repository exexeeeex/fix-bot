import { Navigation } from '@/widgets/navigation';
import { isUserRole } from '@/entities/user';
import { Outlet, redirect } from 'react-router';

export async function clientLoader() {
  if (window.sessionStorage.getItem('_tokens') === null) {
    return redirect('/');
  }

  const role = window.sessionStorage.getItem('_user-role');
  if (!isUserRole(role)) return redirect('/');
  if (role !== 'student') return redirect('/lessons');

  return null;
}

export function HydrateFallback() {
  return <main className="min-h-screen bg-background" aria-busy="true" />;
}

export default function AppLayout() {
  return (
    <>
      <Outlet />
      <Navigation />
    </>
  );
}
