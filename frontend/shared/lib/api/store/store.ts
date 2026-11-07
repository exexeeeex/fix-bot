import { configureStore } from '@reduxjs/toolkit';
import { useDispatch, useSelector } from 'react-redux';

import { baseApi } from './base-api';
import authenticationReducer from '@/features/authentication/slice/authentication';
import applicationReducer from '@/features/application/slices/application';

export const store = configureStore({
  reducer: {
    [baseApi.reducerPath]: baseApi.reducer,
    authentication: authenticationReducer,
    applications: applicationReducer,
  },

  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: ['files', 'registrationThunk'],
        ignoredPaths: ['files', 'registration'],
      },
    }).concat(baseApi.middleware),
});

import './initialize-api';

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export const useAppDispatch = () => useDispatch<AppDispatch>();
export const useAppSelector = useSelector.withTypes<RootState>();
