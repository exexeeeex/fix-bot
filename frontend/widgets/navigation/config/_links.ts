import type { UserRole } from "@/entities/user";
import type { IconType } from "@/shared/lib";

export type Link = {
    id: number;
    name: string;
    href: string;
    icon: IconType;
    access: UserRole | UserRole[];
}

export const LINKS: Link[] = [
    {
        id: 1,
        name: 'Пропуск',
        href: '/',
        icon: 'qr',
        access: 'user'
    },
    {
        id: 2,
        name: 'Профиль',
        href: '/profile',
        icon: 'user',
        access: ['user', 'teacher']
    },
    {
        id: 3,
        name: 'Занятия',
        href: '/lessons',
        icon: 'book',
        access: ['teacher', 'admin']
    },
    {
        id: 4,
        name: 'Статистика',
        href: '/stats',
        icon: 'charts',
        access: ['teacher', 'admin']
    }
] 