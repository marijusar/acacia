import { cache } from 'react';
import { env } from '@/lib/config/env';
import { logger } from '@/lib/config/logger';
import {
  type CreateTeamLLMAPIKeyInput,
  type TeamLLMAPIKeyStatusResponse,
  type TeamLLMAPIKeysListResponse,
  teamLLMAPIKeyStatusResponse,
  teamLLMAPIKeysListResponse,
} from '@/lib/schemas/team-llm-api-keys';
import { BaseHttpService, type BaseServiceArguments } from './base-service';

type TeamLLMAPIKeysServiceArguments = BaseServiceArguments & {};

class TeamLLMAPIKeysService extends BaseHttpService {
  constructor(args: TeamLLMAPIKeysServiceArguments) {
    super(args);
  }

  getAPIKeys = cache(
    async (teamId: number): Promise<TeamLLMAPIKeysListResponse> => {
      const response = await fetch(`${this.url}/teams/${teamId}/llm-api-keys`, {
        method: 'GET',
        headers: { Cookie: await this.cookieService.getAuthCookies() },
      });

      if (!response.ok) {
        const body = await response.json();
        this.logger.error('[GET_API_KEYS] Failed to get API keys', body);
        throw new Error(body.message);
      }

      const unsafeBody = await response.json();
      const body = teamLLMAPIKeysListResponse.parse(unsafeBody);

      return body;
    }
  );

  async createOrUpdateAPIKey(
    teamId: number,
    input: CreateTeamLLMAPIKeyInput
  ): Promise<TeamLLMAPIKeyStatusResponse> {
    const response = await fetch(`${this.url}/teams/${teamId}/llm-api-keys`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Cookie: await this.cookieService.getAuthCookies(),
      },
      body: JSON.stringify(input),
    });

    if (!response.ok) {
      const body = await response.json();
      this.logger.error(
        '[CREATE_OR_UPDATE_API_KEY] Failed to create/update API key',
        body
      );
      throw new Error(body.message);
    }

    const unsafeBody = await response.json();
    const body = teamLLMAPIKeyStatusResponse.parse(unsafeBody);

    return body;
  }

  async deleteAPIKey(teamId: number, keyId: number): Promise<void> {
    const response = await fetch(
      `${this.url}/teams/${teamId}/llm-api-keys/${keyId}`,
      {
        method: 'DELETE',
        headers: {
          Cookie: await this.cookieService.getAuthCookies(),
        },
      }
    );

    if (!response.ok) {
      const body = await response.json();
      this.logger.error('[DELETE_API_KEY] Failed to delete API key', body);
      throw new Error(body.message);
    }
  }
}

export const teamLLMAPIKeysService = new TeamLLMAPIKeysService({
  url: env.ACACIA_API_URL,
  logger,
});
