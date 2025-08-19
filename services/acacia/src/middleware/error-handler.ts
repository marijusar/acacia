import type { Request, Response, NextFunction } from 'express';
import { ZodError } from 'zod';
import { AppError } from '#dashboard-api/errors/app-error.ts';

export const errorHandler = (
  error: unknown,
  request: Request,
  response: Response,
  // Express error handling middleware needs to accept 4 args.
  //eslint-disable-next-line @typescript-eslint/no-unused-vars
  _next: NextFunction
): void => {
  request.log.error(error, 'Request error');

  if (error instanceof ZodError) {
    const validationErrors = error.issues.map((validationError) => ({
      message: validationError.message,
    }));

    response.status(400).json({
      error: 'Validation failed',
      details: validationErrors,
    });
    return;
  }

  if (error instanceof AppError) {
    response.status(error.statusCode).json({
      error: error.message,
    });
    return;
  }

  response.status(500).json({
    error: 'Internal server error',
  });
};
