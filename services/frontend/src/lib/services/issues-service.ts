import { cache } from 'react';
import { env } from '@/lib/config/env';
import { logger } from '@/lib/config/logger';
import { BaseHttpService, BaseServiceArguments } from './base-service';

type CreateIssueParams = {
  name: string;
  column_id: number;
  description?: string;
};

type IssuesServiceArguments = {} & BaseServiceArguments;

class IssuesService extends BaseHttpService {
  constructor(args: IssuesServiceArguments) {
    super(args);
  }

  getProjectIssues = cache(async (projectId: string) => {
    const response = await fetch(`${this.url}/issues/project/${projectId}`, {
      method: 'GET',
    });

    if (!response.ok) {
      const body = await response.json();
      this.logger.error(
        '[GET_PROJECT_ISSUES] Could not get project issues',
        body
      );
      throw new Error('Could not get project issues.');
    }

    const body = await response.json();
    return body;
  });

  async createIssue(params: CreateIssueParams) {
    const response = await fetch(`${this.url}/issues`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(params),
    });

    if (!response.ok) {
      const body = await response.json();
      const message = `[CREATE_ISSUE] Failed to create issue.`;
      this.logger.error(message, body);
      throw new Error(message);
    }

    const body = await response.json();

    return body;
  }
}

export const issuesService = new IssuesService({
  url: env.ACACIA_API_URL,
  logger,
});
