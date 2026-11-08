import { type RouteConfig, index, layout, route } from '@react-router/dev/routes';

export default [
  index('routes/initial.tsx'),
  layout('routes/app-layout.tsx', [
    route('home', 'routes/home.tsx'),
    route('profile', 'routes/test.tsx'),
  ]),
  layout('routes/admin-layout.tsx', [
    route('admin/profile', 'routes/admin-profile.tsx'),
    route('lessons', 'routes/lessons.tsx'),
    route('stats', 'routes/stats.tsx'),
    route('applications', 'routes/applications.tsx'),
  ]),
] satisfies RouteConfig;
