import type { UserRole } from '@/entities/user';
import type { IconType } from '@/shared/lib';

export type Link = {
  id: number;
  name: string;
  href: string;
  icon: IconType;
  access: UserRole | UserRole[];
};

export const LINKS: Link[] = [
  {
    id: 1,
    name: 'Пропуск',
    href: '/home',
    icon: 'qr',
    access: 'student',
  },
  {
    id: 3,
    name: 'Занятия',
    href: '/lessons',
    icon: 'book',
    access: ['teacher', 'admin'],
  },
  {
    id: 4,
    name: 'Статистика',
    href: '/stats',
    icon: 'charts',
    access: ['teacher', 'admin'],
  },
  {
    id: 5,
    name: 'Сканер',
    href: '/admin/profile',
    icon: 'qr',
    access: ['admin', 'teacher'],
  },
  {
    id: 6,
    name: 'Заявки',
    href: '/applications',
    icon: 'user',
    access: ['admin', 'teacher'],
  },
];
