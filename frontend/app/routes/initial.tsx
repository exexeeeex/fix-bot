import { redirect } from 'react-router';
import { isUserRole } from '@/entities/user';
import { Initial } from '@/widgets/initial/ui';

export async function clientLoader() {
  if (window.sessionStorage.getItem('_tokens') === null) return null;

  const role = window.sessionStorage.getItem('_user-role');
  if (!isUserRole(role)) return null;

  return redirect(role === 'student' ? '/home' : '/lessons');
}

export function HydrateFallback() {
  return <main className="min-h-screen bg-background" aria-busy="true" />;
}

export default function InitialPage() {
  return <Initial />;
}
