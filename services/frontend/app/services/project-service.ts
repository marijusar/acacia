import { env } from '~/config/env';
import { logger, type BaseLogger } from '~/config/logger';
import { projectDetailsResponse } from '~/schemas/projects';

class ProjectService {
  private url;
  private logger;
  constructor(url: string, logger: BaseLogger) {
    this.url = url;
    this.logger = logger;
  }
  async getProjectDetails(id: string) {
    console.log(this.url);
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
    body.columns.sort((a, b) => a.position_index - b.position_index);

    return body;
  }
}

export const projectService = new ProjectService(env.ACACIA_API_URL, logger);
