export type UserRole = 'student' | 'admin' | 'teacher';

export const isUserRole = (role: string | null): role is UserRole =>
  role === 'student' || role === 'admin' || role === 'teacher';
