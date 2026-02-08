package pkg

import (
	"time"

	"github.com/cthulhu-platform/common/pkg/env"
)

const (
	BODY_LIMIT_MB = 500

	// Lifecycle TTL for buckets: anonymous vs authorized users
	// LifecycleTTLAnonymous  = 48 * time.Hour
	LifecycleTTLAnonymous  = 5 * time.Minute
	LifecycleTTLAuthorized = 14 * 24 * time.Hour
)

var (
	CORS_ORIGIN = env.GetEnv("CORS_ORIGIN", "*")

	APP_HOST = env.GetEnv("APP_HOST", "0.0.0.0")
	APP_PORT = env.GetEnv("APP_PORT", "7777")

	APP_TEST_ENV = env.GetEnv("APP_TEST_ENV", "")

	// gRPC service URLs
	AUTH_GRPC_URL        = env.GetEnv("AUTH_GRPC_URL", "localhost:49051")
	FILEMANAGER_GRPC_URL = env.GetEnv("FILEMANAGER_GRPC_URL", "localhost:48051")
	LIFECYCLE_GRPC_URL   = env.GetEnv("LIFECYCLE_GRPC_URL", "localhost:50051")
)
