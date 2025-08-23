import express from 'express';
import environment from '#acacia/config/environment.ts';
import { pinoHttp } from 'pino-http';
import { errorHandler } from '#acacia/middleware/error-handler.ts';
import issuesRoutes from '#acacia/routes/issues.ts';
import { databaseInjectionMiddleware } from '#acacia/middleware/database-injection.ts';

const createServer = () => {
  const { port } = environment;
  const app = express();

  app.use(pinoHttp());
  app.use(express.json());
  app.use(databaseInjectionMiddleware);

  app.get('/health', (_request, response) => {
    _request.log.warn('log');
    response.send('Ok nicesu!');
  });

  app.use('/api/issues', issuesRoutes);

  app.use(errorHandler);

  app.listen(port);
};

export { createServer };
