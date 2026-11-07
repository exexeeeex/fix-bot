import type { Middleware } from '@reduxjs/toolkit';
import { allModules } from '.';

export const getApiMiddlewares = (): Middleware[] => {
  const middlewares: Middleware[] = [];
  const added = new Set<string>();

  for (const mod of Object.values(allModules)) {
    for (const api of mod.apis) {
      if (added.has(api.reducerPath)) continue;
      added.add(api.reducerPath);
      middlewares.push(api.middleware);
    }
  }

  return middlewares;
};
