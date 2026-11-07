import type { Application } from './application';

export interface UpdateApplicationStatusRequest {
  id: string;
  status: Application['status'];
}
