import type { Application } from '@/entities/application';
import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import { applicationApi } from '@/entities/application';

export interface ApplicationSliceState {
  applicationsByOrganizationId: Application[] | null;
  applicationsByGroupId: Application[] | null;
  applications: Application[] | null;
  isLoading: boolean;
  error: string | null;
}

const initialState: ApplicationSliceState = {
  applicationsByOrganizationId: null,
  applicationsByGroupId: null,
  applications: null,
  isLoading: false,
  error: null,
};

const applicationsSlice = createSlice({
  name: 'applications',
  initialState,
  reducers: {
    setApplicationsByOrganizationId: (state, action: PayloadAction<Application[] | null>) => {
      state.applicationsByOrganizationId = action.payload;
    },
    setApplicationsByGroupId: (state, action: PayloadAction<Application[] | null>) => {
      state.applicationsByGroupId = action.payload;
    },
    setApplications: (state, action: PayloadAction<Application[] | null>) => {
      state.applications = action.payload;
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.isLoading = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder
      .addMatcher(applicationApi.endpoints.create.matchPending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addMatcher(applicationApi.endpoints.create.matchFulfilled, (state, action) => {
        state.isLoading = false;
        state.error = null;
        state.applications = [
          ...(state.applications ?? []).filter(({ id }) => id !== action.payload.id),
          action.payload,
        ];
      })
      .addMatcher(applicationApi.endpoints.create.matchRejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message ?? 'Не удалось создать заявку';
      })
      .addMatcher(applicationApi.endpoints.getByOrganizationId.matchPending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addMatcher(applicationApi.endpoints.getByOrganizationId.matchFulfilled, (state, action) => {
        state.isLoading = false;
        state.error = null;
        state.applicationsByOrganizationId = action.payload;
        state.applications = mergeApplications(state.applications, action.payload);
      })
      .addMatcher(applicationApi.endpoints.getByOrganizationId.matchRejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message ?? 'Не удалось загрузить заявки организации';
      })
      .addMatcher(applicationApi.endpoints.getByGroupId.matchPending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addMatcher(applicationApi.endpoints.getByGroupId.matchFulfilled, (state, action) => {
        state.isLoading = false;
        state.error = null;
        state.applicationsByGroupId = action.payload;
        state.applications = mergeApplications(state.applications, action.payload);
      })
      .addMatcher(applicationApi.endpoints.getByGroupId.matchRejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message ?? 'Не удалось загрузить заявки группы';
      });
  },
});

const mergeApplications = (current: Application[] | null, incoming: Application[]) => {
  const applicationsById = new Map(
    (current ?? []).map((application) => [application.id, application])
  );
  for (const application of incoming) {
    applicationsById.set(application.id, application);
  }
  return Array.from(applicationsById.values());
};

export const {
  setApplicationsByOrganizationId,
  setApplicationsByGroupId,
  setApplications,
  setError,
  setLoading,
} = applicationsSlice.actions;

export default applicationsSlice.reducer;
