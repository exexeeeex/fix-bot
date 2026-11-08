import { isUserRole } from '@/entities/user';
import { Navigation } from '@/widgets/navigation';
import { Outlet, redirect } from 'react-router';

export async function clientLoader() {
  if (window.sessionStorage.getItem('_tokens') === null) {
    return redirect('/');
  }

  const role = window.sessionStorage.getItem('_user-role');
  if (!isUserRole(role)) return redirect('/');
  if (role === 'student') return redirect('/home');

  return null;
}

export function HydrateFallback() {
  return <main className="min-h-screen bg-background" aria-busy="true" />;
}

export default function AdminLayout() {
  return (
    <>
      <Outlet />
      <Navigation />
    </>
  );
}
