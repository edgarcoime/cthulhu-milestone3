package pkg

import (
	"time"

	"github.com/cthulhu-platform/common/pkg/env"
)

const (
	DEFAULT_CLEANUP_INTERVAL = 6 * time.Minute
)

var (
	APP_HOST = env.GetEnv("APP_HOST", "0.0.0.0")
	APP_PORT = env.GetEnv("APP_PORT", "50051")

	SQLITE_DB_FILE = env.GetEnv("SQLITE_DB_FILE", "lifecycle.db")

	FILEMANAGER_GRPC_URL = env.GetEnv("FILEMANAGER_GRPC_URL", "localhost:48051")
)
