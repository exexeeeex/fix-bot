import type { Group } from '@/entities/group';
import type { BaseEntity } from '@/shared/lib/types';

export interface Application extends BaseEntity {
  id: string;
  name: string;
  status: 'pending' | 'reject' | 'rejected' | 'approved';
  type: 'student' | 'teacher';
  userId: string;
  organizationId: string;
  groupId: string;
  group: Group | null;
}
