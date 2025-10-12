import { cookies } from 'next/headers';

class CookieService {
  /**
   * Captures authentication cookies (access-token and refresh-token) from the incoming request
   * and returns them as a Cookie header string to forward to the backend API
   */
  async getAuthCookies(): Promise<string> {
    const cookieStore = await cookies();

    const accessToken = cookieStore.get('access-token');
    const refreshToken = cookieStore.get('refresh-token');

    const authCookies = [];
    if (accessToken) {
      authCookies.push(`${accessToken.name}=${accessToken.value}`);
    }
    if (refreshToken) {
      authCookies.push(`${refreshToken.name}=${refreshToken.value}`);
    }

    return authCookies.join('; ');
  }

  /**
   * Forwards all Set-Cookie headers from the backend response to the client
   */
  async forwardResponseCookies(response: Response): Promise<void> {
    const setCookieHeaders = response.headers.getSetCookie();
    if (!setCookieHeaders || setCookieHeaders.length === 0) {
      return;
    }

    const cookieStore = await cookies();

    // Forward each Set-Cookie header to the client
    for (const setCookieHeader of setCookieHeaders) {
      // Parse the cookie name=value part (before the first semicolon)
      const [nameValue] = setCookieHeader.split(';');
      const [name, ...valueParts] = nameValue.split('=');
      const value = valueParts.join('='); // Handle values that contain '='

      if (name && value !== undefined) {
        // Simply set the cookie - the browser will handle the attributes from the original Set-Cookie header
        cookieStore.set(name.trim(), value.trim());
      }
    }
  }
}

export const cookieService = new CookieService();
