import { Request, Response } from 'express';
import { database } from '#dashboard-api/config/database.ts';
import { IssuesModel } from '#dashboard-api/models/issues.ts';
import { createIssueSchema, updateIssueSchema, issueParametersSchema } from '#dashboard-api/schemas/issues.ts';

export const getAllIssues = async (request: Request, response: Response): Promise<void> => {
  try {
    const issuesModel = new IssuesModel(database);
    const issues = await issuesModel.findAll();
    response.json(issues);
  } catch (error) {
    request.log.error(error, 'Failed to fetch issues');
    response.status(500).json({ error: 'Internal server error' });
  }
};

export const getIssueById = async (request: Request, response: Response): Promise<void> => {
  try {
    const { id } = issueParametersSchema.parse(request.params);
    const issuesModel = new IssuesModel(database);
    const issue = await issuesModel.findById(id);

    if (!issue) {
      response.status(404).json({ error: 'Issue not found' });
      return;
    }

    response.json(issue);
  } catch (error) {
    if (error instanceof Error && 'issues' in error) {
      response.status(400).json({ error: 'Invalid issue ID' });
      return;
    }
    request.log.error(error, 'Failed to fetch issue');
    response.status(500).json({ error: 'Internal server error' });
  }
};

export const createIssue = async (request: Request, response: Response): Promise<void> => {
  try {
    const data = createIssueSchema.parse(request.body);
    const issuesModel = new IssuesModel(database);
    const issue = await issuesModel.create(data);
    response.status(201).json(issue);
  } catch (error) {
    if (error instanceof Error && 'issues' in error) {
      response.status(400).json({ error: error.message });
      return;
    }
    request.log.error(error, 'Failed to create issue');
    response.status(500).json({ error: 'Internal server error' });
  }
};

export const updateIssue = async (request: Request, response: Response): Promise<void> => {
  try {
    const { id } = issueParametersSchema.parse(request.params);
    const data = updateIssueSchema.parse(request.body);

    if (Object.keys(data).length === 0) {
      response.status(400).json({ error: 'No fields provided for update' });
      return;
    }

    const issuesModel = new IssuesModel(database);
    const issue = await issuesModel.update(id, data);

    if (!issue) {
      response.status(404).json({ error: 'Issue not found' });
      return;
    }

    response.json(issue);
  } catch (error) {
    if (error instanceof Error && 'issues' in error) {
      response.status(400).json({ error: error.message });
      return;
    }
    request.log.error(error, 'Failed to update issue');
    response.status(500).json({ error: 'Internal server error' });
  }
};

export const deleteIssue = async (request: Request, response: Response): Promise<void> => {
  try {
    const { id } = issueParametersSchema.parse(request.params);
    const issuesModel = new IssuesModel(database);
    const deletedIssue = await issuesModel.delete(id);

    if (!deletedIssue) {
      response.status(404).json({ error: 'Issue not found' });
      return;
    }

    response.status(204).send();
  } catch (error) {
    if (error instanceof Error && 'issues' in error) {
      response.status(400).json({ error: 'Invalid issue ID' });
      return;
    }
    request.log.error(error, 'Failed to delete issue');
    response.status(500).json({ error: 'Internal server error' });
  }
};