import { baseApi } from '@/shared/lib/api/store/base-api';
import { isSessionTokens, storeSessionTokens } from '@/shared/lib/api/session-storage';
import type { AuthenticationResponse } from '../types';

export function storeAuthenticatedSession(data: AuthenticationResponse) {
  if (typeof window === 'undefined') return;
  if (!data.success) throw new Error('Authentication was not successful');

  storeSessionTokens(data.tokens);
  window.sessionStorage.setItem('_user-role', data.user.role);
  if (data.user.organization_id) {
    window.sessionStorage.setItem('_organization-id', data.user.organization_id);
  } else {
    window.sessionStorage.removeItem('_organization-id');
  }
}

export const authenticationApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    login: builder.mutation<AuthenticationResponse, { initData: string }>({
      query: ({ initData }) => ({
        url: 'auth/max',
        method: 'POST',
        body: { initData },
      }),
    }),
    demoLogin: builder.mutation<AuthenticationResponse, { role: 'student' | 'teacher' }>({
      query: (body) => ({
        url: 'auth/demo',
        method: 'POST',
        body,
      }),
    }),
    refreshTokens: builder.mutation<AuthenticationResponse['tokens'], { refresh_token: string }>({
      query: (body) => ({
        url: 'auth/max/refresh',
        method: 'POST',
        body,
      }),
      transformResponse: (response: unknown) => {
        if (!isSessionTokens(response)) {
          throw new Error('Invalid refresh token response');
        }
        return response;
      },
    }),
  }),
});

export const { useLoginMutation, useDemoLoginMutation, useRefreshTokensMutation } =
  authenticationApi;
