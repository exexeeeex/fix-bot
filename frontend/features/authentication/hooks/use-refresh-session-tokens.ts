import { useEffect } from 'react';
import {
  clearSessionStorage,
  getSessionTokens,
  storeSessionTokens,
} from '@/shared/lib/api/session-storage';
import { useRefreshTokensMutation } from '../api/authentication';

export const useRefreshSessionTokens = () => {
  const [refreshTokens] = useRefreshTokensMutation();

  useEffect(() => {
    const tokens = getSessionTokens();
    if (!tokens) return;

    void refreshTokens({ refresh_token: tokens.refresh_token })
      .unwrap()
      .then(storeSessionTokens)
      .catch((error: unknown) => {
        if (
          typeof error === 'object' &&
          error !== null &&
          'status' in error &&
          (error.status === 401 || error.status === 403)
        ) {
          clearSessionStorage();
          return;
        }

        console.error('Failed to refresh session tokens', error);
      });
  }, [refreshTokens]);
};
