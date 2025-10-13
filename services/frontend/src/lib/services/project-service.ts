import { cache } from 'react';
import { env } from '@/lib/config/env';
import { logger } from '@/lib/config/logger';
import {
  projectDetailsResponse,
  projectsResponse,
  projectColumnsResponse,
  createProjectResponse,
  Issue,
} from '@/lib/schemas/projects';
import { BaseHttpService, BaseServiceArguments } from './base-service';

type CreateProjectParams = {
  name: string;
  team_id: number;
};

type CreateProjectColumnParams = {
  name: string;
  project_id: number;
};

type ProjectServiceArguments = BaseServiceArguments & {};

class ProjectService extends BaseHttpService {
  constructor(args: ProjectServiceArguments) {
    super(args);
  }

  getProjectDetails = cache(async (id: string) => {
    const authCookies = await this.cookieService.getAuthCookies();

    const response = await fetch(`${this.url}/projects/${id}/details`, {
      method: 'GET',
      headers: {
        Cookie: authCookies,
      },
    });

    if (!response.ok) {
      const body = await response.json();
      this.logger.error(
        '[GET_PROJECT_DETAILS] Could not get project details',
        body
      );
      return undefined;
    }

    const unsafeBody = await response.json();

    const body = projectDetailsResponse.parse(unsafeBody);
    // sort so columns are in order.
    body.columns.sort((a, b) => a.position_index - b.position_index);

    const issueMap: Record<number, Issue[]> = body.issues.reduce(
      (acc: Record<number, Issue[]>, issue) => {
        if (!acc[issue.column_id]) {
          acc[issue.column_id] = [];
        }
        acc[issue.column_id].push(issue);

        return acc;
      },
      {}
    );

    return {
      ...body,
      issues: issueMap,
    };
  });

  getProjects = cache(async () => {
    const authCookies = await this.cookieService.getAuthCookies();

    const response = await fetch(`${this.url}/projects/`, {
      method: 'GET',
      headers: {
        Cookie: authCookies,
      },
    });

    if (!response.ok) {
      const body = await response.json();
      this.logger.error('[GET_PROJECTS] Could not get projects.', body);
      return undefined;
    }

    const unsafeBody = await response.json();

    const body = projectsResponse.parse(unsafeBody);

    return body;
  });

  getProjectColumns = cache(async (projectId: string) => {
    const authCookies = await this.cookieService.getAuthCookies();

    const response = await fetch(
      `${this.url}/project-columns/project/${projectId}`,
      {
        headers: {
          Cookie: authCookies,
        },
      }
    );

    if (!response.ok) {
      const body = await response.json();
      const message = `[GET_PROJECT_COLUMNS] Failed to get project columns.`;
      this.logger.error(message, body);
      throw new Error(message);
    }

    const body = await response.json();
    return projectColumnsResponse.parse(body);
  });

  // Mutation methods (not cached)
  async createProject(params: CreateProjectParams) {
    const authCookies = await this.cookieService.getAuthCookies();

    const response = await fetch(`${this.url}/projects`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Cookie: authCookies,
      },
      body: JSON.stringify(params),
    });

    if (!response.ok) {
      const body = await response.json();
      const message = `[CREATE_PROJECT] Failed to create project.`;
      this.logger.error(message, body);
      throw new Error(message);
    }

    const body = await response.json();
    return createProjectResponse.parse(body);
  }

  async createProjectColumn(params: CreateProjectColumnParams) {
    const authCookies = await this.cookieService.getAuthCookies();

    const response = await fetch(`${this.url}/project-columns/`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Cookie: authCookies,
      },
      body: JSON.stringify(params),
    });

    if (!response.ok) {
      const body = await response.json();
      const message = `[CREATE_PROJECT_COLUMN] Failed to create project column`;
      this.logger.error(message, body);
      throw new Error(message);
    }

    return true;
  }

  async deleteProjectColumn(id: number) {
    const authCookies = await this.cookieService.getAuthCookies();

    const response = await fetch(`${this.url}/project-columns/${id}`, {
      method: 'DELETE',
      headers: {
        Cookie: authCookies,
      },
    });

    if (!response.ok) {
      const body = await response.json();
      const message = `[DELETE_PROJECT_COLUMN] Failed to delete project column`;
      this.logger.error(message, body);
      throw new Error(message);
    }

    return true;
  }
}

export const projectService = new ProjectService({
  url: env.ACACIA_API_URL,
  logger,
});
