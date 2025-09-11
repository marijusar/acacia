import { type RouteConfig, index, route } from '@react-router/dev/routes';

export default [
  index('routes/home.tsx'),
  route('projects/board', 'projects/board/board.tsx'),
] satisfies RouteConfig;
