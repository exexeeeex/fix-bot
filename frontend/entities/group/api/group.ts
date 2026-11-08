import { baseApi } from '@/shared/lib/api';
import type { Group } from '../types';

export const groupApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    listByOrganization: builder.query<Group[], string>({
      query: (organizationId) => `groups/by-org/${organizationId}`,
    }),
  }),
});

export const { useListByOrganizationQuery: useGroupsByOrganizationQuery } = groupApi;
