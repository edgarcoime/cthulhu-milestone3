package pkg

import (
	"time"

	"github.com/cthulhu-platform/common/pkg/env"
)

const (
	DEFAULT_REPOSITORY_QUERY_TIMEOUT = 5 * time.Second
	BUCKET_TOKEN_EXPIRATION          = 30 * time.Minute
	PRESIGNED_URL_EXPIRATION         = 15 * time.Minute
)

var (
	APP_HOST      = env.GetEnv("APP_HOST", "0.0.0.0")
	APP_PORT      = env.GetEnv("APP_PORT", "48051")
	AUTH_GRPC_URL = env.GetEnv("AUTH_GRPC_URL", "localhost:49051")

	SQLITE_DB_FILE          = env.GetEnv("SQLITE_DB_FILE", "filemanager.db")
	BUCKET_TOKEN_SECRET_KEY = env.GetEnv("BUCKET_TOKEN_SECRET_KEY", "iamasecretkey")

	S3_ACCESS_KEY_ID        = env.GetEnv("S3_ACCESS_KEY_ID", "")
	S3_SECRET_ACCESS_KEY    = env.GetEnv("S3_SECRET_ACCESS_KEY", "")
	S3_ENDPOINT             = env.GetEnv("S3_ENDPOINT", "")
	S3_PRESIGNED_ENDPOINT   = env.GetEnv("S3_PRESIGNED_ENDPOINT", "") // if set, presigned URLs use this (e.g. localhost:4566 for browser); server-side still uses S3_ENDPOINT
	S3_REGION               = env.GetEnv("S3_REGION", "")
	S3_BUCKET_NAME          = env.GetEnv("S3_BUCKET_NAME", "")
	S3_FORCE_PATH_STYLE     = env.GetEnv("S3_FORCE_PATH_STYLE", "true")
)
