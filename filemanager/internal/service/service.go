// Package service implements the filemanager business logic. Uploads use a
// two-phase presigned URL flow: PrepareUpload returns presigned PUT URLs;
// the client uploads directly to S3; ConfirmUpload persists file metadata.
package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cthulhu-platform/filemanager/internal/connections"
	"github.com/cthulhu-platform/filemanager/internal/repository"
	"github.com/cthulhu-platform/filemanager/internal/storage"
	"github.com/cthulhu-platform/filemanager/pkg"
	pb "github.com/cthulhu-platform/proto/pkg/filemanager"
)

// Service is the filemanager application service.
type Service interface {
	// Upload (presigned URL flow): PrepareUpload then client PUTs to S3 then ConfirmUpload.
	PrepareUpload(ctx context.Context, req *pb.PrepareUploadRequest) (*pb.PrepareUploadResponse, error)
	ConfirmUpload(ctx context.Context, req *pb.ConfirmUploadRequest) (*pb.ConfirmUploadResponse, error)

	// Download (presigned GET URL; for protected buckets, bucket_access_token required)
	PrepareDownload(ctx context.Context, req *pb.PrepareDownloadRequest) (*pb.PrepareDownloadResponse, error)

	// Bucket metadata and auth
	RetrieveFileBucket(ctx context.Context, storageID string) (*pkg.BucketMetadata, error)
	GetBucketAdmins(ctx context.Context, bucketID string) (*pkg.BucketAdminsResponse, error)
	IsBucketProtected(ctx context.Context, bucketID string) (bool, *string, error)
	AuthenticateBucket(ctx context.Context, bucketID string, password string, userID *string, authTokenID *string) (string, error)

	// DeleteBucket removes the bucket, its files (and S3 objects), and bucket_admins. Returns files deleted count.
	DeleteBucket(ctx context.Context, bucketID string) (filesDeleted int64, err error)
}

type filemanagerService struct {
	repo    repository.Repository
	storage storage.Storage
	conns   *connections.ConnectionsContainer
}

func NewFilemanagerService(repo repository.Repository, stor storage.Storage, conns *connections.ConnectionsContainer) Service {
	return &filemanagerService{
		repo:    repo,
		storage: stor,
		conns:   conns,
	}
}

func (s *filemanagerService) RetrieveFileBucket(ctx context.Context, storageID string) (*pkg.BucketMetadata, error) {
	_, err := s.repo.GetBucketByID(ctx, storageID)
	if err != nil {
		return nil, err
	}
	files, err := s.repo.GetFilesByBucketID(ctx, storageID)
	if err != nil {
		return nil, err
	}
	out := &pkg.BucketMetadata{StorageID: storageID, Files: make([]pkg.FileInfo, 0, len(files))}
	var totalSize int64
	for _, f := range files {
		out.Files = append(out.Files, pkg.FileInfo{
			OriginalName: f.OriginalName,
			StringID:     f.StringID,
			Key:          f.S3Key,
			Size:         f.Size,
			ContentType:  f.ContentType,
		})
		totalSize += f.Size
	}
	out.TotalSize = totalSize
	return out, nil
}

func (s *filemanagerService) GetBucketAdmins(ctx context.Context, bucketID string) (*pkg.BucketAdminsResponse, error) {
	_, err := s.repo.GetBucketByID(ctx, bucketID)
	if err != nil {
		return nil, err
	}
	list, err := s.repo.GetBucketAdminsByBucketID(ctx, bucketID)
	if err != nil {
		return nil, err
	}
	out := &pkg.BucketAdminsResponse{BucketID: bucketID, Admins: make([]pkg.AdminInfo, 0, len(list))}
	for i, a := range list {
		info := pkg.AdminInfo{
			UserID:    a.UserID,
			Email:     "",
			Username:  nil,
			AvatarURL: nil,
			IsOwner:   i == 0,
			CreatedAt: a.CreatedAt,
		}
		out.Admins = append(out.Admins, info)
	}
	if len(out.Admins) > 0 {
		out.Owner = &out.Admins[0]
	}
	return out, nil
}

func (s *filemanagerService) IsBucketProtected(ctx context.Context, bucketID string) (bool, *string, error) {
	bucket, err := s.repo.GetBucketByID(ctx, bucketID)
	if err != nil {
		return false, nil, err
	}
	return bucket.PasswordHash.Valid, nil, nil
}

func (s *filemanagerService) AuthenticateBucket(ctx context.Context, bucketID string, password string, userID *string, authTokenID *string) (string, error) {
	bucket, err := s.repo.GetBucketByID(ctx, bucketID)
	if err != nil {
		return "", err
	}
	if !bucket.PasswordHash.Valid {
		return "", fmt.Errorf("bucket is not protected")
	}
	if !VerifyBucketPassword(password, bucket.PasswordHash.String) {
		return "", fmt.Errorf("invalid password")
	}
	return GenerateBucketAccessToken(bucketID, userID, authTokenID, []string{"read"})
}

func (s *filemanagerService) DeleteBucket(ctx context.Context, bucketID string) (filesDeleted int64, err error) {
	_, err = s.repo.GetBucketByID(ctx, bucketID)
	if err != nil {
		return 0, fmt.Errorf("bucket not found: %w", err)
	}
	files, err := s.repo.GetFilesByBucketID(ctx, bucketID)
	if err != nil {
		return 0, fmt.Errorf("list files: %w", err)
	}
	for _, f := range files {
		if delErr := s.storage.DeleteObject(ctx, f.S3Key); delErr != nil {
			slog.Warn("failed to delete S3 object during bucket delete", "s3_key", f.S3Key, "error", delErr)
		}
	}
	if err := s.repo.DeleteBucket(ctx, bucketID); err != nil {
		return 0, fmt.Errorf("delete bucket: %w", err)
	}
	return int64(len(files)), nil
}
