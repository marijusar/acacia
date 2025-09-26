import {
  type RouteConfig,
  index,
  layout,
  route,
} from '@react-router/dev/routes';

export default [
  index('routes/home.tsx'),
  route('/api/projects/create', 'api/projects/create.ts'),
  route('/api/columns/create', 'api/columns/create.ts'),
  route('/api/columns/remove', 'api/columns/remove.ts'),

  layout('./layouts/dashboard.tsx', [
    route('projects/board/:id', 'projects/board/board.tsx'),
    route(
      'projects/board/:id/settings/columns',
      'projects/board/column-settings.tsx'
    ),
  ]),
] satisfies RouteConfig;
