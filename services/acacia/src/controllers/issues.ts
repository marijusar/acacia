import type { Request, Response, NextFunction } from 'express';
import { database } from '#dashboard-api/config/database.ts';
import { IssuesModel } from '#dashboard-api/models/issues.ts';
import { AppError } from '#dashboard-api/errors/app-error.ts';
import {
  createIssueSchema,
  updateIssueSchema,
  issueParametersSchema,
} from '#dashboard-api/schemas/issues.ts';

export const getAllIssues = async (
  _request: Request,
  response: Response,
  next: NextFunction
): Promise<void> => {
  try {
    const issuesModel = new IssuesModel(database);
    const issues = await issuesModel.findAll();
    response.json(issues);
  } catch (error) {
    next(error);
  }
};

export const getIssueById = async (
  request: Request,
  response: Response,
  next: NextFunction
): Promise<void> => {
  try {
    const { id } = issueParametersSchema.parse(request.params);
    const issuesModel = new IssuesModel(database);
    const issue = await issuesModel.findById(id);

    if (!issue) {
      throw new AppError('Issue not found', 404);
    }

    response.json(issue);
  } catch (error) {
    next(error);
  }
};

export const createIssue = async (
  request: Request,
  response: Response,
  next: NextFunction
): Promise<void> => {
  try {
    const data = createIssueSchema.parse(request.body);
    const issuesModel = new IssuesModel(database);
    const issue = await issuesModel.create(data);
    response.status(201).json(issue);
  } catch (error) {
    next(error);
  }
};

export const updateIssue = async (
  request: Request,
  response: Response,
  next: NextFunction
): Promise<void> => {
  try {
    const { id } = issueParametersSchema.parse(request.params);
    const data = updateIssueSchema.parse(request.body);

    if (Object.keys(data).length === 0) {
      throw new AppError('No fields provided for update', 400);
    }

    const issuesModel = new IssuesModel(database);
    const issue = await issuesModel.update(id, data);

    if (!issue) {
      throw new AppError('Issue not found', 404);
    }

    response.json(issue);
  } catch (error) {
    next(error);
  }
};

export const deleteIssue = async (
  request: Request,
  response: Response,
  next: NextFunction
): Promise<void> => {
  try {
    const { id } = issueParametersSchema.parse(request.params);
    const issuesModel = new IssuesModel(database);
    const deletedIssue = await issuesModel.delete(id);

    if (!deletedIssue) {
      throw new AppError('Issue not found', 404);
    }

    response.status(204).send();
  } catch (error) {
    next(error);
  }
};
