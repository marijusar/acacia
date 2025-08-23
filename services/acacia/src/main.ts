import environment from '#dashboard-api/config/environment.ts';
import express from 'express';
import { pinoHttp } from 'pino-http';
import issuesRoutes from '#dashboard-api/routes/issues.ts';
import { errorHandler } from '#dashboard-api/middleware/error-handler.ts';
