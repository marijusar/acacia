import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

const ACACIA_API_URL = process.env.ACACIA_API_URL;

if (!ACACIA_API_URL) {
  throw new Error('ACACIA_API_URL environment variable is not set');
}

export async function middleware(request: NextRequest) {
  const accessToken = request.cookies.get('access-token');
  const refreshToken = request.cookies.get('refresh-token');

  // If user has both tokens, proceed with the request
  if (accessToken && refreshToken) {
    return NextResponse.next();
  }

  // If user has no tokens at all, redirect to login
  if (!accessToken && !refreshToken) {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  // If user only has refresh token, call auth endpoint to get new access token
  if (!accessToken && refreshToken) {
    try {
      const response = await fetch(`${ACACIA_API_URL}/users/auth/me`, {
        method: 'GET',
        headers: {
          Cookie: `refresh-token=${refreshToken.value}`,
        },
      });

      // If unauthorized, redirect to login
      if (!response.ok && response.status === 401) {
        return NextResponse.redirect(new URL('/login', request.url));
      }

      // If successful, attach the refreshed cookies to the response
      if (response.ok) {
        const nextResponse = NextResponse.next();

        // Forward Set-Cookie headers from backend
        const setCookieHeaders = response.headers.getSetCookie();
        for (const setCookieHeader of setCookieHeaders) {
          nextResponse.headers.append('Set-Cookie', setCookieHeader);
        }

        return nextResponse;
      }

      // For any other error, redirect to login
      return NextResponse.redirect(new URL('/login', request.url));
    } catch (error) {
      // On network error or other issues, redirect to login
      console.error('Auth middleware error:', error);
      return NextResponse.redirect(new URL('/login', request.url));
    }
  }

  // Fallback: redirect to login
  return NextResponse.redirect(new URL('/login', request.url));
}

export const config = {
  matcher: [
    /*
     * Match only /projects and /teams routes
     * Static assets are automatically excluded by Next.js
     */
    '/projects/:path*',
    '/teams/:path*',
  ],
};
