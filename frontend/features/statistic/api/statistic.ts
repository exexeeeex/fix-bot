import { type Abcense } from '@/entities/abcense';
import type { Student } from '@/entities/student';
import { baseApi } from '@/shared/lib/api';

interface AbcenseResponse {
  all_students: Student[];
  abcenses_students: Student[];
  late_students: Student[];
}

const mapApplicationResponse = (abcense: AbcenseResponse): Abcense => ({
  allStudents: abcense.all_students,
  abcensesStudents: abcense.abcenses_students,
  lateStudents: abcense.late_students ?? [],
});

export const statisticApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getAbcenses: builder.query<Abcense, string>({
      query: (organizationId) => ({
        url: `statistic/organization-abcenses/${organizationId}`,
        method: 'GET',
      }),
      transformResponse: mapApplicationResponse,
    }),
  }),
});

export const { useGetAbcensesQuery } = statisticApi;
