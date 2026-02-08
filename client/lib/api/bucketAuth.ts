import { API_URL } from '@/lib/config';
import { tokenStorage, ensureValidToken } from './userAuth';

export const bucketTokenStorage = {
  getBucketAccessToken: (bucketId: string): string | null => {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem(`bucket_access_${bucketId}`);
  },

  setBucketAccessToken: (bucketId: string, token: string): void => {
    if (typeof window === 'undefined') return;
    localStorage.setItem(`bucket_access_${bucketId}`, token);
  },

  clearBucketAccessToken: (bucketId: string): void => {
    if (typeof window === 'undefined') return;
    localStorage.removeItem(`bucket_access_${bucketId}`);
  },
};

export const checkBucketProtected = async (
  bucketId: string
): Promise<{ protected: boolean; bucket_id: string }> => {
  const response = await fetch(`${API_URL}/files/s/${bucketId}/protected`);

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Failed to check bucket protection' }));
    const err = new Error(error.error || 'Failed to check bucket protection') as Error & { status?: number };
    err.status = response.status;
    throw err;
  }

  return response.json();
};

export const authenticateBucket = async (
  bucketId: string,
  password: string
): Promise<{ access_token: string; expires_in: number }> => {
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
  };

  const accessToken = tokenStorage.getAccessToken();
  if (accessToken) {
    try {
      const validToken = await ensureValidToken();
      (headers as Record<string, string>)['Authorization'] = `Bearer ${validToken}`;
    } catch {
      // Continue without auth token if validation fails
    }
  }

  const response = await fetch(`${API_URL}/files/s/${bucketId}/authenticate`, {
    method: 'POST',
    headers,
    body: JSON.stringify({ password }),
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Authentication failed' }));
    throw new Error(error.error || 'Authentication failed');
  }

  return response.json();
};
