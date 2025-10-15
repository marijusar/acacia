import { redirect } from 'next/navigation';
import { cache } from 'react';
import { env } from '@/lib/config/env';
import { logger } from '@/lib/config/logger';
import {
  authStatusResponse,
  type AuthStatusResponse,
  loginResponse,
  type LoginResponse,
  type LoginInput,
  registerResponse,
  type RegisterResponse,
  type RegisterInput,
} from '@/lib/schemas/users';
import { BaseHttpService, type BaseServiceArguments } from './base-service';

type UserServiceArguments = BaseServiceArguments & {};

class UserService extends BaseHttpService {
  constructor(args: UserServiceArguments) {
    super(args);
  }

  async login(credentials: LoginInput): Promise<LoginResponse> {
    const response = await fetch(`${this.url}/users/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(credentials),
    });

    // Forward Set-Cookie headers (access-token and refresh-token) to client
    await this.cookieService.forwardResponseCookies(response);

    if (!response.ok) {
      const body = await response.json();
      this.logger.error('[LOGIN] Login failed', body);
      throw new Error(body.message);
    }

    const unsafeBody = await response.json();
    const body = loginResponse.parse(unsafeBody);

    return body;
  }

  async register(data: RegisterInput): Promise<RegisterResponse> {
    const response = await fetch(`${this.url}/users/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const body = await response.json();
      this.logger.error('[REGISTER] Registration failed', body);
      throw new Error(body.message);
    }

    const unsafeBody = await response.json();
    const body = registerResponse.parse(unsafeBody);

    return body;
  }
}

export const userService = new UserService({
  url: env.ACACIA_API_URL,
  logger,
});
