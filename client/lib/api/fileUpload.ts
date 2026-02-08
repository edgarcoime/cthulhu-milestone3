import { API_URL } from '@/lib/config';
import { tokenStorage, ensureValidToken } from './userAuth';
import type { FileInfo } from './types';

export type { FileInfo };

export interface PrepareUploadFileMeta {
  original_name: string;
  size: number;
  content_type: string;
}

export interface PrepareUploadRequest {
  files: PrepareUploadFileMeta[];
  password?: string;
}

export interface UploadSlot {
  string_id: string;
  presigned_put_url: string;
  s3_key: string;
}

export interface PrepareUploadResponse {
  storage_id: string;
  slots: UploadSlot[];
  error?: string;
}

export interface ConfirmUploadFileMeta {
  string_id: string;
  original_name: string;
  size: number;
  content_type: string;
}

export interface ConfirmUploadRequest {
  storage_id: string;
  files: ConfirmUploadFileMeta[];
}

export interface ConfirmUploadResponse {
  success: boolean;
  storage_id?: string;
  files?: FileInfo[];
  total_size?: number;
  error?: string;
}

function buildUploadHeaders(): Promise<HeadersInit> {
  const headers: HeadersInit = { 'Content-Type': 'application/json' };
  const accessToken = tokenStorage.getAccessToken();
  if (!accessToken) return Promise.resolve(headers);
  return ensureValidToken()
    .then((token) => {
      (headers as Record<string, string>)['Authorization'] = `Bearer ${token}`;
      return headers;
    })
    .catch(() => {
      throw new Error('Authentication failed - please sign in again');
    });
}

export const uploadFile = async (
  files: FileList,
  password?: string
): Promise<{ status: boolean; data?: { url: string; storage_id: string }; error?: string }> => {
  if (files.length === 0) {
    throw new Error('No files to upload');
  }

  const filesMeta: PrepareUploadFileMeta[] = [];
  for (let i = 0; i < files.length; i++) {
    const f = files[i];
    filesMeta.push({
      original_name: f.name,
      size: f.size,
      content_type: f.type || 'application/octet-stream',
    });
  }

  const prepareBody: PrepareUploadRequest = { files: filesMeta };
  if (password && password.trim()) {
    prepareBody.password = password.trim();
  }

  const headers = await buildUploadHeaders();

  const prepareRes = await fetch(`${API_URL}/files/upload/prepare`, {
    method: 'POST',
    headers,
    body: JSON.stringify(prepareBody),
  });

  const prepareData: PrepareUploadResponse = await prepareRes.json();
  if (!prepareRes.ok || prepareData.error) {
    throw new Error(prepareData.error || 'Prepare upload failed');
  }
  if (!prepareData.storage_id || !prepareData.slots?.length || prepareData.slots.length !== files.length) {
    throw new Error('Invalid prepare response');
  }

  const { storage_id, slots } = prepareData;

  for (let i = 0; i < slots.length; i++) {
    const file = files[i];
    const slot = slots[i];
    const putRes = await fetch(slot.presigned_put_url, {
      method: 'PUT',
      headers: {
        'Content-Type': file.type || 'application/octet-stream',
      },
      body: file,
    });
    if (!putRes.ok) {
      throw new Error(`Upload failed for ${file.name}`);
    }
  }

  const confirmBody: ConfirmUploadRequest = {
    storage_id,
    files: slots.map((slot, i) => ({
      string_id: slot.string_id,
      original_name: files[i].name,
      size: files[i].size,
      content_type: files[i].type || 'application/octet-stream',
    })),
  };

  const confirmRes = await fetch(`${API_URL}/files/upload/confirm`, {
    method: 'POST',
    headers,
    body: JSON.stringify(confirmBody),
  });

  const confirmData: ConfirmUploadResponse = await confirmRes.json();
  if (!confirmRes.ok || !confirmData.success || confirmData.error) {
    throw new Error(confirmData.error || 'Confirm upload failed');
  }

  const url = confirmData.storage_id ? `/files/s/${confirmData.storage_id}` : '';

  return {
    status: true,
    data: {
      url,
      storage_id: confirmData.storage_id || storage_id,
    },
  };
};
