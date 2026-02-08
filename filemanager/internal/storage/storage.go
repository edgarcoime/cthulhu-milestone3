package storage

import "context"

type Storage interface {
	Close() error
	// PresignPut returns a short-lived presigned URL for uploading an object via PUT.
	PresignPut(ctx context.Context, key string, contentLength int64, contentType string) (url string, err error)
	// PresignGet returns a short-lived presigned URL for downloading an object via GET.
	PresignGet(ctx context.Context, key string) (url string, err error)
	// DeleteObject deletes an object from storage by key. NoSuchKey is treated as success.
	DeleteObject(ctx context.Context, key string) error
}
