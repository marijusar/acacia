import environment from '#dashboard-api/config/environment.ts';
import express from 'express';
import { pinoHttp } from 'pino-http';
import issuesRoutes from '#dashboard-api/routes/issues.ts';

const { port } = environment;
const app = express();

app.use(pinoHttp());
app.use(express.json());

app.get('/health', (_request, response) => {
  _request.log.warn('log');
  response.send('Ok nicesu!');
});

app.use('/api/issues', issuesRoutes);

app.listen(port);
