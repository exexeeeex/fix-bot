import { baseApi } from '@/shared/lib/api';
import type { Group } from '@/entities/group';
import {
  type Application,
  type CreateApplicationRequest,
  type UpdateApplicationStatusRequest,
} from '../types';

interface ApplicationResponse {
  id: string;
  name: string;
  status: Application['status'];
  type: Application['type'];
  user_id: string;
  organization_id: string;
  group_id: string | null;
  group?: Group | null;
  created_at: string;
  updated_at: string;
}

const mapApplicationResponse = (application: ApplicationResponse): Application => ({
  id: application.id,
  name: application.name,
  status: application.status,
  type: application.type,
  userId: application.user_id,
  organizationId: application.organization_id,
  groupId: application.group_id ?? '',
  group: application.group ?? null,
  createdAt: application.created_at,
  updatedAt: application.updated_at,
});

export const applicationApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    create: builder.mutation<Application, CreateApplicationRequest>({
      query: ({ name, organizationId, groupId, userId, type }) => ({
        url: 'application/send',
        method: 'POST',
        body: {
          name,
          organization_id: organizationId,
          group_id: type === 'teacher' ? null : groupId || null,
          user_id: userId,
          type,
        },
      }),
      transformResponse: mapApplicationResponse,
    }),
    getByOrganizationId: builder.query<Application[], string>({
      query: (id) => ({
        url: `application/by-organization/${id}`,
        method: 'GET',
      }),
      transformResponse: (response: ApplicationResponse[] | null) =>
        response?.map(mapApplicationResponse) ?? [],
      providesTags: (_result, _error, id) => [{ type: 'Applications', id: `organization:${id}` }],
    }),
    getByGroupId: builder.query<Application[], string>({
      query: (id) => ({
        url: `application/by-group/${id}`,
        method: 'GET',
      }),
      transformResponse: (response: ApplicationResponse[] | null) =>
        response?.map(mapApplicationResponse) ?? [],
      providesTags: (_result, _error, id) => [{ type: 'Applications', id: `group:${id}` }],
    }),
    updateStatus: builder.mutation<Application, UpdateApplicationStatusRequest>({
      query: ({ id, status }) => ({
        url: 'application/update-status',
        method: 'PUT',
        body: { id, status },
      }),
      transformResponse: mapApplicationResponse,
      invalidatesTags: (application) =>
        application
          ? [
              { type: 'Applications', id: `organization:${application.organizationId}` },
              ...(application.groupId
                ? [{ type: 'Applications' as const, id: `group:${application.groupId}` }]
                : []),
            ]
          : [],
    }),
  }),
});

export const {
  useCreateMutation,
  useGetByOrganizationIdQuery,
  useGetByGroupIdQuery,
  useUpdateStatusMutation,
} = applicationApi;
