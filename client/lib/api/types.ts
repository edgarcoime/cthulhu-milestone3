export interface FileInfo {
  original_name: string;
  string_id: string;
  key: string;
  size: number;
  content_type: string;
}

export interface BucketMetadata {
  storage_id: string;
  files: FileInfo[];
  total_size: number;
}

export interface AdminInfo {
  user_id: string;
  email: string;
  username?: string;
  avatar_url?: string;
  is_owner: boolean;
  created_at: number;
}

export interface BucketAdminsResponse {
  bucket_id: string;
  owner: AdminInfo | null;
  admins: AdminInfo[];
}

export interface UploadFileResponse {
  transaction_id: string;
  success: boolean;
  error?: string;
  storage_id?: string;
  files?: FileInfo[];
  total_size?: number;
}
