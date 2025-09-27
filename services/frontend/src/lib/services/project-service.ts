import { cache } from 'react';
import { env } from '@/lib/config/env';
import { logger, type BaseLogger } from '@/lib/config/logger';
import {
  projectDetailsResponse,
  projectsResponse,
  projectColumnsResponse,
  createProjectResponse,
} from '@/lib/schemas/projects';

type CreateProjectParams = {
  name: string;
};

type CreateProjectColumnParams = {
  name: string;
  project_id: number;
};

class ProjectService {
  private url;
  private logger;

  constructor(url: string, logger: BaseLogger) {
    this.url = url;
    this.logger = logger;
  }

  getProjectDetails = cache(async (id: string) => {
    const response = await fetch(`${this.url}/projects/${id}/details`, {
      method: 'GET',
    });

    if (!response.ok) {
      const body = await response.json();
      this.logger.error(
        '[GET_PROJECT_DETAILS] Could not get project details',
        body
      );
      throw new Error('Could not get project details.');
    }

    const unsafeBody = await response.json();

    const body = projectDetailsResponse.parse(unsafeBody);
    // sort so columns are in order.
    body.columns.sort((a, b) => a.position_index - b.position_index);

    return body;
  });

  getProjects = cache(async () => {
    const response = await fetch(`${this.url}/projects/`, {
      method: 'GET',
    });

    if (!response.ok) {
      const body = await response.json();
      this.logger.error('[GET_PROJECTS] Could not get projects.', body);
      throw new Error('Could not get projects');
    }

    const unsafeBody = await response.json();

    const body = projectsResponse.parse(unsafeBody);

    return body;
  });

  getProjectColumns = cache(async (projectId: string) => {
    const response = await fetch(
      `${this.url}/project-columns/project/${projectId}`
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
    const response = await fetch(`${this.url}/projects`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
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
    const response = await fetch(`${this.url}/project-columns/`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
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
    const response = await fetch(`${this.url}/project-columns/${id}`, {
      method: 'DELETE',
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

export const projectService = new ProjectService(env.ACACIA_API_URL, logger);

