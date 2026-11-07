import type { ModuleConfig } from '.';
import {
  applicationApi,
  authenticationApi,
  groupApi,
  organizationApi,
  qrApi,
} from './initialize-api';

export const authenticationModule: ModuleConfig = {
  apis: [authenticationApi],
};

export const applicationModule: ModuleConfig = {
  key: 'application',
  apis: [applicationApi],
  injectOn: ['/'],
};

export const organizationModule: ModuleConfig = {
  key: 'organization',
  apis: [organizationApi],
  injectOn: ['/'],
};

export const groupModule: ModuleConfig = {
  key: 'group',
  apis: [groupApi],
  injectOn: ['/'],
};

export const qrModule: ModuleConfig = {
  key: 'qr',
  apis: [qrApi],
  injectOn: ['/home'],
};

export const allModules: Record<string, ModuleConfig> = {
  authentication: authenticationModule,
  application: applicationModule,
  organization: organizationModule,
  group: groupModule,
  qr: qrModule,
};

export const coreModules = Object.entries(allModules)
  .filter(([_, module]) => module.injectImmediately)
  .map(([name, _]) => name);

export const lazyModules = Object.entries(allModules)
  .filter(([_, module]) => !module.injectImmediately)
  .map(([name, _]) => name);
