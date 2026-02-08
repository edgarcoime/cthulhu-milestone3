import { API_URL } from '@/lib/config';

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: {
    id: string;
    email: string;
    username?: string;
    avatar_url?: string;
  };
}

export interface TokenPair {
  access_token: string;
  refresh_token: string;
}

export interface Claims {
  user_id: string;
  email: string;
  provider: string;
}

export const tokenStorage = {
  getAccessToken: (): string | null => {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('access_token');
  },

  getRefreshToken: (): string | null => {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('refresh_token');
  },

  setTokens: (accessToken: string, refreshToken: string): void => {
    if (typeof window === 'undefined') return;
    localStorage.setItem('access_token', accessToken);
    localStorage.setItem('refresh_token', refreshToken);
    window.dispatchEvent(new Event('auth-state-changed'));
  },

  clearTokens: (): void => {
    if (typeof window === 'undefined') return;
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    window.dispatchEvent(new Event('auth-state-changed'));
  },
};

export const isAuthenticated = (): boolean => {
  if (typeof window === 'undefined') return false;
  return !!tokenStorage.getAccessToken();
};

export const initiateOAuth = (provider: string = 'github'): void => {
  if (typeof window === 'undefined') return;
  window.location.href = `${API_URL}/auth/oauth/${provider}`;
};

export const handleCallback = async (
  code: string,
  state: string,
  provider: string = 'github'
): Promise<AuthResponse> => {
  const response = await fetch(
    `${API_URL}/auth/oauth/${provider}/callback?code=${encodeURIComponent(code)}&state=${encodeURIComponent(state)}`,
    {
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
    }
  );

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Failed to authenticate' }));
    throw new Error(error.error || 'Failed to authenticate');
  }

  const data: AuthResponse = await response.json();
  tokenStorage.setTokens(data.access_token, data.refresh_token);
  return data;
};

export const validateToken = async (token: string): Promise<Claims> => {
  const response = await fetch(`${API_URL}/auth/validate`, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Invalid token' }));
    throw new Error(error.error || 'Invalid token');
  }

  const data = await response.json();
  return data.claims;
};

export const refreshToken = async (refreshTokenValue: string): Promise<TokenPair> => {
  const response = await fetch(`${API_URL}/auth/refresh`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ refresh_token: refreshTokenValue }),
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Failed to refresh token' }));
    throw new Error(error.error || 'Failed to refresh token');
  }

  const data: TokenPair = await response.json();
  tokenStorage.setTokens(data.access_token, data.refresh_token);
  return data;
};

export const logout = async (refreshTokenValue: string): Promise<void> => {
  try {
    await fetch(`${API_URL}/auth/logout`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ refresh_token: refreshTokenValue }),
    });
  } finally {
    tokenStorage.clearTokens();
  }
};

export const hardLogout = (): void => {
  tokenStorage.clearTokens();
  if (typeof window !== 'undefined') {
    window.location.href = '/signin';
  }
};

export const ensureValidToken = async (): Promise<string> => {
  const accessToken = tokenStorage.getAccessToken();
  const refreshTokenValue = tokenStorage.getRefreshToken();

  if (!accessToken || !refreshTokenValue) {
    throw new Error('No authentication tokens available');
  }

  try {
    await validateToken(accessToken);
    return accessToken;
  } catch {
    try {
      const newTokens = await refreshToken(refreshTokenValue);
      return newTokens.access_token;
    } catch {
      hardLogout();
      throw new Error('Token refresh failed - user logged out');
    }
  }
};

export const getCurrentUserId = async (): Promise<string | null> => {
  const accessToken = tokenStorage.getAccessToken();
  if (!accessToken) return null;

  try {
    const claims = await validateToken(accessToken);
    return claims.user_id;
  } catch {
    return null;
  }
};
