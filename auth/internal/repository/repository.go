package repository

import (
	"context"

	"github.com/cthulhu-platform/auth/internal/repository/sqlc/db"
)

type Repository interface {
	Close() error

	// User operations
	GetUserByOAuthID(ctx context.Context, provider string, userID string) (*db.User, error)
	GetUserByID(ctx context.Context, id string) (*db.User, error)
	GetUserByEmail(ctx context.Context, email string) (*db.User, error)
	CreateUser(ctx context.Context, user *db.User) error
	UpdateUser(ctx context.Context, user *db.User) error
	SoftDeleteUser(ctx context.Context, id string) error

	// Refresh token operations
	CreateRefreshToken(ctx context.Context, token *db.RefreshToken) error
	GetRefreshTokenByHash(ctx context.Context, hash string) (*db.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, id string, reason string) error
	RevokeAllUserTokens(ctx context.Context, userID string, reason string) error

	// OAuth session operations
	CreateOAuthSession(ctx context.Context, session *db.OauthSession) error
	GetOAuthSession(ctx context.Context, state string) (*db.OauthSession, error)
	DeleteOAuthSession(ctx context.Context, state string) error
}
