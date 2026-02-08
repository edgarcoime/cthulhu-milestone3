import { API_URL } from '@/lib/config';
import { bucketTokenStorage } from './bucketAuth';
import { getCurrentUserId } from './userAuth';
import type { BucketMetadata, BucketAdminsResponse } from './types';

export type { BucketMetadata, BucketAdminsResponse };
export type { AdminInfo } from './types';

export interface BucketLifecycleResponse {
  bucket_id: string;
  expires_at: string; // UTC ISO 8601
}

export const fetchBucketLifecycle = async (
  bucketId: string
): Promise<BucketLifecycleResponse> => {
  const headers: HeadersInit = {};
  const bucketToken = bucketTokenStorage.getBucketAccessToken(bucketId);
  if (bucketToken) {
    (headers as Record<string, string>)['X-Bucket-Token'] = bucketToken;
  }

  const response = await fetch(`${API_URL}/lifecycle/s/${bucketId}`, { headers });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Lifecycle not found' }));
    throw new Error(error.error || 'Lifecycle not found');
  }

  return response.json();
};

export const fetchBucketFiles = async (bucketId: string): Promise<BucketMetadata> => {
  const headers: HeadersInit = {};
  const bucketToken = bucketTokenStorage.getBucketAccessToken(bucketId);
  if (bucketToken) {
    (headers as Record<string, string>)['X-Bucket-Token'] = bucketToken;
  }

  const response = await fetch(`${API_URL}/files/s/${bucketId}`, { headers });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Failed to fetch files' }));
    throw new Error(error.error || 'Failed to fetch files');
  }

  return response.json();
};

/**
 * @param bucketId - Bucket (storage) ID
 * @param stringId - File string_id (URL segment, may be UUID-like)
 * @param downloadAs - Filename for the saved file (e.g. original_name from file list). Always prefer passing this so the download uses the original name; if omitted, stringId is used.
 */
export const downloadFile = async (
  bucketId: string,
  stringId: string,
  downloadAs?: string
): Promise<void> => {
  const bucketToken = bucketTokenStorage.getBucketAccessToken(bucketId);
  const url = `${API_URL}/files/s/${bucketId}/d/${stringId}`;

  const headers: HeadersInit = {};
  if (bucketToken) {
    (headers as Record<string, string>)['X-Bucket-Token'] = bucketToken;
  }

  const response = await fetch(url, { headers });

  if (response.ok) {
    const blob = await response.blob();
    const downloadUrl = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = downloadUrl;
    const name = (downloadAs?.trim() ?? '') || stringId;
    a.download = name;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(downloadUrl);
    return;
  }

  // Protected bucket without token: fall back to new tab (filename may be wrong)
  if (response.status === 401 && !bucketToken) {
    window.open(url, '_blank');
    return;
  }

  throw new Error('Download failed');
};

export const fetchBucketAdmins = async (bucketId: string): Promise<BucketAdminsResponse> => {
  const headers: HeadersInit = {};
  const bucketToken = bucketTokenStorage.getBucketAccessToken(bucketId);
  if (bucketToken) {
    (headers as Record<string, string>)['X-Bucket-Token'] = bucketToken;
  }

  const response = await fetch(`${API_URL}/files/s/${bucketId}/admins`, { headers });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Failed to fetch admins' }));
    throw new Error(error.error || 'Failed to fetch admins');
  }

  return response.json();
};

export const isBucketAdmin = async (bucketId: string): Promise<boolean> => {
  try {
    const userId = await getCurrentUserId();
    if (!userId) return false;

    const admins = await fetchBucketAdmins(bucketId);
    return (
      admins.admins.some((admin) => admin.user_id === userId) ||
      (admins.owner !== null && admins.owner.user_id === userId)
    );
  } catch {
    return false;
  }
};
