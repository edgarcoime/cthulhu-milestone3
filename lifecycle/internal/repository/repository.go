package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"os"
	"path/filepath"
	"time"

	commonrepo "github.com/cthulhu-platform/common/pkg/repository"
	internalpkg "github.com/cthulhu-platform/lifecycle/internal/pkg"
	"github.com/cthulhu-platform/lifecycle/pkg"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schemas/up.sql
var upSQL string

type Repository interface {
	Close() error
	PutLifecycle(ctx context.Context, lifecycle pkg.Lifecycle) (*pkg.Lifecycle, error)
	GetLifecycle(ctx context.Context, bucketSlug string) (*pkg.Lifecycle, error)
	DeleteLifecycle(ctx context.Context, bucketSlug string) error
	ListExpiredLifecycles(ctx context.Context, now time.Time) ([]pkg.Lifecycle, error)
}

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

	// Open SQlite database connection
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Printf("Failed to open SQLite database: %v\n", err)
		db.Close()
		return nil, err
	}

	// Ping database to ensure connection is established
	if err := db.PingContext(ctx); err != nil {
		log.Printf("Failed to ping SQLite database: %v\n", err)
		db.Close()
		return nil, err
	}

	// Initialize database schema
	if _, err := db.ExecContext(ctx, upSQL); err != nil {
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

func (r *sqliteRepository) PutLifecycle(ctx context.Context, lifecycle pkg.Lifecycle) (*pkg.Lifecycle, error) {
	query := `
		INSERT INTO lifecycle (bucket_slug, expires_at)
		VALUES ($1, $2)
		ON CONFLICT (bucket_slug) DO UPDATE SET
			expires_at = excluded.expires_at,
			updated_at = datetime('now')
	`
	_, err := r.db.ExecContext(ctx, query, lifecycle.BucketSlug, lifecycle.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return r.GetLifecycle(ctx, lifecycle.BucketSlug)
}

func (r *sqliteRepository) GetLifecycle(ctx context.Context, bucketSlug string) (*pkg.Lifecycle, error) {
	query := `
		SELECT id, bucket_slug, expires_at, created_at, updated_at
		FROM lifecycle
		WHERE bucket_slug = $1
	`
	row := r.db.QueryRowContext(ctx, query, bucketSlug)
	lifecycle := &pkg.Lifecycle{}
	var expiresAt, createdAt, updatedAt string
	err := row.Scan(&lifecycle.ID, &lifecycle.BucketSlug, &expiresAt, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	lifecycle.ExpiresAt, err = commonrepo.ParseSQLiteTime(expiresAt)
	if err != nil {
		return nil, err
	}
	lifecycle.CreatedAt, err = commonrepo.ParseSQLiteTime(createdAt)
	if err != nil {
		return nil, err
	}
	lifecycle.UpdatedAt, err = commonrepo.ParseSQLiteTime(updatedAt)
	if err != nil {
		return nil, err
	}
	return lifecycle, nil
}

func (r *sqliteRepository) DeleteLifecycle(ctx context.Context, bucketSlug string) error {
	query := `
		DELETE FROM lifecycle
		WHERE bucket_slug = $1
	`
	_, err := r.db.ExecContext(ctx, query, bucketSlug)
	if err != nil {
		return err
	}
	return nil
}

func (r *sqliteRepository) ListExpiredLifecycles(ctx context.Context, now time.Time) ([]pkg.Lifecycle, error) {
	query := `
		SELECT id, bucket_slug, expires_at, created_at, updated_at
		FROM lifecycle
		WHERE expires_at < $1
	`
	rows, err := r.db.QueryContext(ctx, query, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []pkg.Lifecycle
	for rows.Next() {
		var l pkg.Lifecycle
		var expiresAt, createdAt, updatedAt string
		if err := rows.Scan(&l.ID, &l.BucketSlug, &expiresAt, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		var parseErr error
		l.ExpiresAt, parseErr = commonrepo.ParseSQLiteTime(expiresAt)
		if parseErr != nil {
			return nil, parseErr
		}
		l.CreatedAt, parseErr = commonrepo.ParseSQLiteTime(createdAt)
		if parseErr != nil {
			return nil, parseErr
		}
		l.UpdatedAt, parseErr = commonrepo.ParseSQLiteTime(updatedAt)
		if parseErr != nil {
			return nil, parseErr
		}
		out = append(out, l)
	}
	return out, rows.Err()
}
