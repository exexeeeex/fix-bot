import { createApi } from '@reduxjs/toolkit/query/react';
import baseQueryWithReauth from '../base-query';

export const baseApi = createApi({
  reducerPath: 'api',
  baseQuery: baseQueryWithReauth,
  tagTypes: ['Applications'],
  endpoints: () => ({}),
});
