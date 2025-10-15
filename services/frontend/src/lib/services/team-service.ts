import { cache } from 'react';
import { env } from '@/lib/config/env';
import { logger } from '@/lib/config/logger';
import {
  userTeamsResponse,
  type UserTeamsResponse,
  createTeamResponse,
  type CreateTeamResponse,
  type CreateTeamInput,
} from '@/lib/schemas/teams';
import { BaseHttpService, type BaseServiceArguments } from './base-service';

type TeamServiceArguments = BaseServiceArguments & {};

class TeamService extends BaseHttpService {
  constructor(args: TeamServiceArguments) {
    super(args);
  }

  getUserTeams = cache(async (): Promise<UserTeamsResponse> => {
    const authCookies = await this.cookieService.getAuthCookies();

    const response = await fetch(`${this.url}/teams/`, {
      method: 'GET',
      headers: {
        Cookie: authCookies,
      },
    });

    if (!response.ok) {
      const body = await response.json();
      this.logger.error('[GET_USER_TEAMS] Could not get user teams', body);
      throw new Error('Could not get user teams');
    }

    const unsafeBody = await response.json();
    const body = userTeamsResponse.parse(unsafeBody);

    return body;
  });

  async createTeam(params: CreateTeamInput): Promise<CreateTeamResponse> {
    const authCookies = await this.cookieService.getAuthCookies();

    const response = await fetch(`${this.url}/teams`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Cookie: authCookies,
      },
      body: JSON.stringify(params),
    });

    if (!response.ok) {
      const body = await response.json();
      this.logger.error('[CREATE_TEAM] Failed to create team', body);
      throw new Error(body.message);
    }

    const unsafeBody = await response.json();
    return createTeamResponse.parse(unsafeBody);
  }
}

export const teamService = new TeamService({
  url: env.ACACIA_API_URL,
  logger,
});
