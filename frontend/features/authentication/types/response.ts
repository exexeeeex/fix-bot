import type { UserRole } from '@/entities/user';

export interface AuthenticationResponse {
  success: boolean;
  user: {
    id: string;
    first_name: string;
    last_name: string;
    username: string;
    language_code: string | null;
    photo_url: string | null;
    role: UserRole;
    organization_id: string;
    applicationStatus: 'none' | 'pending' | 'approved' | 'rejected';
  };
  tokens: TokensResponse;
}

export interface TokensResponse {
  access_token: string;
  refresh_token: string;
}
