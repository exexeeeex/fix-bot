import { useEffect, useState } from 'react';
import { isUserRole, type UserRole } from '@/entities/user';
import { LINKS } from '../config';
import { NavigationItem } from './item';
import { useAppSelector } from '@/shared/lib/api/store/store';

export const Navigation = () => {
  const user = useAppSelector((state) => state.authentication.user);
  const [storedRole, setStoredRole] = useState<UserRole | null>(null);

  useEffect(() => {
    const role = window.sessionStorage.getItem('_user-role');
    if (isUserRole(role)) setStoredRole(role);
  }, []);

  const role = user?.role ?? storedRole;
  const links = LINKS.filter((link) => {
    if (!role) return false;

    if (Array.isArray(link.access)) {
      return link.access.includes(role);
    }

    return link.access === role;
  });

  return (
    <nav className="fixed bottom-0 left-0 w-full h-15 bg-white dark:bg-gray-950 flex flex-col items-center justify-center">
      <ul className="flex flex-row justify-center items-center gap-5 w-full h-15">
        {links.map((link) => (
          <NavigationItem key={link.id} item={link} />
        ))}
      </ul>
    </nav>
  );
};
