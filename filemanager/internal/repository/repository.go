package repository

import (
	"context"

	"github.com/cthulhu-platform/filemanager/internal/repository/sqlc/db"
)

type Repository interface {
	Close() error

	// Bucket operations
	GetBucketByID(ctx context.Context, id string) (*db.Bucket, error)
	CreateBucket(ctx context.Context, bucket *db.Bucket) error
	UpdateBucket(ctx context.Context, bucket *db.Bucket) error
	DeleteBucket(ctx context.Context, id string) error
	ListBuckets(ctx context.Context, limit int, offset int) ([]*db.Bucket, error)

	// File operations
	GetFileByID(ctx context.Context, id int64) (*db.File, error)
	GetFileByStringID(ctx context.Context, stringID string) (*db.File, error)
	GetFileByBucketIDAndStringID(ctx context.Context, bucketID, stringID string) (*db.File, error)
	GetFilesByBucketID(ctx context.Context, bucketID string) ([]*db.File, error)
	CreateFile(ctx context.Context, file *db.File) error
	UpdateFile(ctx context.Context, file *db.File) error
	DeleteFile(ctx context.Context, id int64) error
	ListFiles(ctx context.Context, limit int, offset int) ([]*db.File, error)

	// Bucket admin operations
	AddBucketAdmin(ctx context.Context, bucketAdmin *db.BucketAdmin) error
	RemoveBucketAdmin(ctx context.Context, userID string, bucketID string) error
	GetBucketAdminsByBucketID(ctx context.Context, bucketID string) ([]*db.BucketAdmin, error)
	GetBucketsByAdminUserID(ctx context.Context, userID string) ([]*db.Bucket, error)
	IsBucketAdmin(ctx context.Context, userID string, bucketID string) (bool, error)
}
