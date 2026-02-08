package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	internalpkg "github.com/cthulhu-platform/auth/internal/pkg"
	"github.com/cthulhu-platform/auth/internal/repository/sqlc/db"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed sqlc/schema.sql
var schemaSQL string

type sqliteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(ctx context.Context) (*sqliteRepository, error) {
	// create timeout context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	path := internalpkg.SQLITE_DB_FILE

	// Ensure directory exists
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	// Open SQLite database connection
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Printf("Failed to open SQLite database: %v\n", err)
		return nil, err
	}

	// Ping database to ensure connection is established
	if err := db.PingContext(ctx); err != nil {
		log.Printf("Failed to ping SQLite database: %v\n", err)
		db.Close()
		return nil, err
	}

	// Initialize database schema (same SQL used by sqlc for code generation)
	if err := runSchema(ctx, db, schemaSQL); err != nil {
		log.Printf("Failed to initialize database schema: %v\n", err)
		db.Close()
		return nil, err
	}

	log.Println("SQLite database initialized successfully")

	r := &sqliteRepository{db: db}
	return r, nil
}

func (r *sqliteRepository) Close() error {
	return r.db.Close()
}

func (r *sqliteRepository) GetUserByOAuthID(ctx context.Context, provider string, userID string) (*db.User, error) {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	user, err := db.New(r.db).GetUserByOAuthID(ctx, db.GetUserByOAuthIDParams{
		OauthProvider: provider,
		OauthUserID:   userID,
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *sqliteRepository) GetUserByID(ctx context.Context, id string) (*db.User, error) {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	user, err := db.New(r.db).GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *sqliteRepository) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	user, err := db.New(r.db).GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *sqliteRepository) CreateUser(ctx context.Context, user *db.User) error {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	err := db.New(r.db).CreateUser(ctx, db.CreateUserParams{
		ID:            user.ID,
		OauthProvider: user.OauthProvider,
		OauthUserID:   user.OauthUserID,
		Email:         user.Email,
		Username:      user.Username,
		AvatarUrl:     user.AvatarUrl,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	})
	return err
}

func (r *sqliteRepository) UpdateUser(ctx context.Context, user *db.User) error {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	return db.New(r.db).UpdateUser(ctx, db.UpdateUserParams{
		ID:        user.ID,
		Username:  user.Username,
		AvatarUrl: user.AvatarUrl,
		UpdatedAt: user.UpdatedAt,
	})
}

func (r *sqliteRepository) SoftDeleteUser(ctx context.Context, id string) error {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	now := time.Now().Unix()
	return db.New(r.db).SoftDeleteUser(ctx, db.SoftDeleteUserParams{
		ID:        id,
		DeletedAt: sql.NullInt64{Int64: now, Valid: true},
		UpdatedAt: now,
	})
}

// Refresh token operations
func (r *sqliteRepository) CreateRefreshToken(ctx context.Context, token *db.RefreshToken) error {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	return db.New(r.db).CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		ID:        token.ID,
		UserID:    token.UserID,
		TokenHash: token.TokenHash,
		ExpiresAt: token.ExpiresAt,
		CreatedAt: token.CreatedAt,
	})
}

func (r *sqliteRepository) GetRefreshTokenByHash(ctx context.Context, hash string) (*db.RefreshToken, error) {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	token, err := db.New(r.db).GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *sqliteRepository) RevokeRefreshToken(ctx context.Context, id string, reason string) error {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	now := time.Now().Unix()
	return db.New(r.db).RevokeRefreshToken(ctx, db.RevokeRefreshTokenParams{
		ID:            id,
		RevokedAt:     sql.NullInt64{Int64: now, Valid: true},
		RevokedReason: sql.NullString{String: reason, Valid: true},
	})
}

func (r *sqliteRepository) RevokeAllUserTokens(ctx context.Context, userID string, reason string) error {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	now := time.Now().Unix()
	return db.New(r.db).RevokeAllUserTokens(ctx, db.RevokeAllUserTokensParams{
		UserID:        userID,
		RevokedAt:     sql.NullInt64{Int64: now, Valid: true},
		RevokedReason: sql.NullString{String: reason, Valid: true},
	})
}

// OAuth session operations

func (r *sqliteRepository) CreateOAuthSession(ctx context.Context, session *db.OauthSession) error {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	return db.New(r.db).CreateOAuthSession(ctx, db.CreateOAuthSessionParams{
		State:         session.State,
		Provider:      session.Provider,
		CodeVerifier:  session.CodeVerifier,
		CodeChallenge: session.CodeChallenge,
		RedirectUri:   session.RedirectUri,
		ExpiresAt:     session.ExpiresAt,
		CreatedAt:     session.CreatedAt,
	})
}

func (r *sqliteRepository) GetOAuthSession(ctx context.Context, state string) (*db.OauthSession, error) {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	session, err := db.New(r.db).GetOAuthSession(ctx, state)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sqliteRepository) DeleteOAuthSession(ctx context.Context, state string) error {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	return db.New(r.db).DeleteOAuthSession(ctx, state)
}

// runSchema executes schema SQL statement by statement (database/sql runs one per Exec).
func runSchema(ctx context.Context, db *sql.DB, schema string) error {
	for _, stmt := range splitStatements(schema) {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}
	return nil
}

func splitStatements(schema string) []string {
	var out []string
	for _, s := range strings.Split(schema, ";") {
		if t := strings.TrimSpace(s); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func defaultTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
