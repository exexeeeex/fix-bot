import { BASE_API_URL } from '@/shared/config';
import {
  fetchBaseQuery,
  type BaseQueryFn,
  type FetchArgs,
  type FetchBaseQueryError,
} from '@reduxjs/toolkit/query';
import {
  clearSessionStorage,
  getSessionTokens,
  isSessionTokens,
  storeSessionTokens,
} from './session-storage';

export const baseQuery = fetchBaseQuery({
  baseUrl: BASE_API_URL,
  credentials: 'include',
  timeout: 15000,
  prepareHeaders: (headers) => {
    const tokens = getSessionTokens();
    if (tokens) {
      headers.set('Authorization', `Bearer ${tokens.access_token}`);
    }
    return headers;
  },
});

let refreshPromise: Promise<boolean> | null = null;

const baseQueryWithReauth: BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError> = async (
  args,
  api,
  extra
) => {
  let result = await baseQuery(args, api, extra);

  if (api.endpoint === 'refreshTokens') return result;

  if (result?.error?.status === 401) {
    if (typeof window === 'undefined') return result;

    if (refreshPromise) {
      if (await refreshPromise) {
        return await baseQuery(args, api, extra);
      }
      return result;
    }

    refreshPromise = (async () => {
      try {
        const tokens = getSessionTokens();
        if (!tokens) {
          clearSessionStorage();
          return false;
        }

        const refreshResult = await baseQuery(
          {
            url: 'auth/max/refresh',
            method: 'POST',
            body: { refresh_token: tokens.refresh_token },
          },
          api,
          extra
        );
        if (refreshResult.error || !isSessionTokens(refreshResult.data)) {
          if (
            !refreshResult.error ||
            refreshResult.error.status === 401 ||
            refreshResult.error.status === 403
          ) {
            clearSessionStorage();
          }
          return false;
        }

        storeSessionTokens(refreshResult.data);
        return true;
      } finally {
        refreshPromise = null;
      }
    })();

    if (await refreshPromise) {
      result = await baseQuery(args, api, extra);
    }
  }

  return result;
};

export default baseQueryWithReauth;
