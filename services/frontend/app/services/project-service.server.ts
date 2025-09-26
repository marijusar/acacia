import { env } from '~/config/env';
import { logger, type BaseLogger } from '~/config/logger';
import {
  createProjectResponse,
  projectColumnsResponse,
  projectDetailsResponse,
  projectsResponse,
} from '~/schemas/projects';

type CreateProjectParams = {
  name: string;
};

type CreateProjectColumnParams = {
  name: string;
  projectId: string;
};

class ProjectService {
  private url;
  private logger;
  constructor(url: string, logger: BaseLogger) {
    this.url = url;
    this.logger = logger;
  }
  async getProjectDetails(id: string) {
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
  }

  async getProjects() {
    const response = await fetch(`${this.url}/projects/`, {
      method: 'GET',
    });

    if (!response.ok) {
      const body = await response.json();
      this.logger.error('[GET_PROJECT_DETAILS] Could not get projects.', body);
      throw new Error('Could not get projects');
    }

    const unsafeBody = await response.json();

    const body = projectsResponse.parse(unsafeBody);
    // sort so columns are in order.

    return body;
  }

  async createProject(params: CreateProjectParams) {
    const response = await fetch(`${this.url}/projects`, {
      method: 'POST',
      body: JSON.stringify(params),
    });

    const body = await response.json();

    if (!response.ok) {
      const body = await response.json();
      const message = `[CREATE_PROJECT] Failed to create project.`;
      this.logger.error(message, body);
      throw new Error(message);
    }

    return createProjectResponse.parse(body);
  }

  async getProjectColumns(projectId: string) {
    const response = await fetch(
      `${this.url}/project-columns/project/${projectId}`
    );

    const body = await response.json();

    if (!response.ok) {
      const body = await response.json();
      const message = `[CREATE_PROJECT] Failed to create project.`;
      this.logger.error(message, body);
      throw new Error(message);
    }

    return projectColumnsResponse.parse(body);
  }

  async createProjectColumn(params: CreateProjectParams) {
    const response = await fetch(`${this.url}/project-columns/`, {
      body: JSON.stringify(params),
      method: 'POST',
    });

    const body = await response.json();

    if (!response.ok) {
      const message = `[CREATE_PROJECT] Failed to create project column`;
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
      const message = `[CREATE_PROJECT] Failed to delete a project column`;
      this.logger.error(message, body);
      throw new Error(message);
    }

    return true;
  }
}

export const projectService = new ProjectService(env.ACACIA_API_URL, logger);
