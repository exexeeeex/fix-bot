import type { Middleware, Reducer } from '@reduxjs/toolkit';

export interface ApiModule {
  reducerPath: string;
  reducer: Reducer;
  middleware: Middleware;
}

export interface StateModule {
  key: string;
  reducer: Reducer;
}

export type Module = ApiModule | StateModule;

export interface InjectableReducer {
  key: string;
  reducer: Reducer;
  replace?: boolean;
}

export interface ModuleConfig {
  key?: string;
  apis: any[];
  injectOn?: string[];
  injectImmediately?: boolean;
}
