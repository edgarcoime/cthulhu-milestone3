// Config
export { API_URL } from '@/lib/config';

// Shared types
export type {
  FileInfo,
  BucketMetadata,
  AdminInfo,
  BucketAdminsResponse,
  UploadFileResponse,
} from './types';

// User auth
export {
  tokenStorage,
  isAuthenticated,
  initiateOAuth,
  handleCallback,
  validateToken,
  refreshToken,
  logout,
  hardLogout,
  ensureValidToken,
  getCurrentUserId,
} from './userAuth';
export type { AuthResponse, TokenPair, Claims } from './userAuth';

// File upload
export {
  uploadFile,
} from './fileUpload';
export type {
  PrepareUploadFileMeta,
  PrepareUploadRequest,
  UploadSlot,
  PrepareUploadResponse,
  ConfirmUploadFileMeta,
  ConfirmUploadRequest,
  ConfirmUploadResponse,
} from './fileUpload';

// Bucket auth (token storage, protected check, authenticate)
export {
  bucketTokenStorage,
  checkBucketProtected,
  authenticateBucket,
} from './bucketAuth';

// Bucket (files, download, admins, lifecycle)
export {
  fetchBucketFiles,
  downloadFile,
  fetchBucketAdmins,
  isBucketAdmin,
  fetchBucketLifecycle,
} from './bucket';
export type { BucketLifecycleResponse } from './bucket';
