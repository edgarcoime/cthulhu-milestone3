package pkg

import (
	"time"

	"github.com/cthulhu-platform/common/pkg/env"
)

const (
	OAUTH_SESSION_EXPIRATION_TIME = 10 * time.Minute
	ACCESS_TOKEN_EXPIRATION       = 15 * time.Minute
	REFRESH_TOKEN_EXPIRATION      = 7 * 24 * time.Hour
)

var (
	APP_HOST = env.GetEnv("APP_HOST", "0.0.0.0")
	APP_PORT = env.GetEnv("APP_PORT", "49051")

	SQLITE_DB_FILE = env.GetEnv("SQLITE_DB_FILE", "auth.db")
	JWT_SECRET_KEY = env.GetEnv("JWT_SECRET", "secret")

	GITHUB_CLIENT_ID     = env.GetEnv("GITHUB_CLIENT_ID", "")
	GITHUB_CLIENT_SECRET = env.GetEnv("GITHUB_CLIENT_SECRET", "")
	GITHUB_REDIRECT_URI  = env.GetEnv("GITHUB_REDIRECT_URI", "")
)
