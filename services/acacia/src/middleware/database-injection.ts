import { createDatabase } from '#acacia/config/database.ts';
import environment from '#acacia/config/environment.ts';
import { NextFunction, Request, Response } from 'express';

const databaseInjectionMiddleware = (
  req: Request,
  _res: Response,
  next: NextFunction
) => {
  const database = createDatabase(environment.database_url);
  req.acacia.db = database;
  next();
};

export { databaseInjectionMiddleware };
