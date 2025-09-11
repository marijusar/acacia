import { type RouteConfig, index, route } from '@react-router/dev/routes';

export default [
  index('routes/home.tsx'),
  route('projects/board/:id', 'projects/board/board.tsx'),
] satisfies RouteConfig;
