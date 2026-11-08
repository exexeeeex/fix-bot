import { baseApi } from '@/shared/lib/api';
import type { Organization } from '../types';

export const organizationApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    list: builder.query<Organization[], void>({
      query: () => 'organizations',
    }),
    searchByName: builder.query<Organization[], string>({
      query: (name) => ({
        url: 'organizations/by-name',
        method: 'POST',
        body: {
          name: name,
        },
      }),
    }),
  }),
});

export const {
  useListQuery: useOrganizationsQuery,
  useSearchByNameQuery: useOrganizationsByNameQuery,
} = organizationApi;
