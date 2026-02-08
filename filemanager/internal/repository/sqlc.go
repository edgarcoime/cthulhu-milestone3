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

	internalpkg "github.com/cthulhu-platform/filemanager/internal/pkg"
	"github.com/cthulhu-platform/filemanager/internal/repository/sqlc/db"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed sqlc/schema.sql
var schemaSQL string

type sqliteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(ctx context.Context) (*sqliteRepository, error) {
	// Create timeout context
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

	// Initialize database schema
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

// Implement Repository interface

// Bucket operations
func (r *sqliteRepository) GetBucketByID(ctx context.Context, id string) (*db.Bucket, error) {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	bucket, err := db.New(r.db).GetBucketByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &bucket, nil
}

func (r *sqliteRepository) CreateBucket(ctx context.Context, bucket *db.Bucket) error {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	return db.New(r.db).CreateBucket(ctx, db.CreateBucketParams{
		ID:           bucket.ID,
		PasswordHash: bucket.PasswordHash,
		CreatedAt:    bucket.CreatedAt,
		UpdatedAt:    bucket.UpdatedAt,
	})
}

func (r *sqliteRepository) UpdateBucket(ctx context.Context, bucket *db.Bucket) error {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	return db.New(r.db).UpdateBucket(ctx, db.UpdateBucketParams{
		ID:           bucket.ID,
		PasswordHash: bucket.PasswordHash,
		UpdatedAt:    bucket.UpdatedAt,
	})
}

func (r *sqliteRepository) DeleteBucket(ctx context.Context, id string) error {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	return db.New(r.db).DeleteBucket(ctx, id)
}

func (r *sqliteRepository) ListBuckets(ctx context.Context, limit int, offset int) ([]*db.Bucket, error) {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	list, err := db.New(r.db).ListBuckets(ctx, db.ListBucketsParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, err
	}
	out := make([]*db.Bucket, 0, len(list))
	for i := range list {
		b := list[i]
		out = append(out, &b)
	}
	return out, nil
}

// File operations
func (r *sqliteRepository) GetFileByID(ctx context.Context, id int64) (*db.File, error) {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	file, err := db.New(r.db).GetFileByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *sqliteRepository) GetFileByStringID(ctx context.Context, stringID string) (*db.File, error) {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	file, err := db.New(r.db).GetFileByStringID(ctx, stringID)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *sqliteRepository) GetFileByBucketIDAndStringID(ctx context.Context, bucketID, stringID string) (*db.File, error) {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	file, err := db.New(r.db).GetFileByBucketIDAndStringID(ctx, db.GetFileByBucketIDAndStringIDParams{
		BucketID: bucketID,
		StringID: stringID,
	})
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *sqliteRepository) GetFilesByBucketID(ctx context.Context, bucketID string) ([]*db.File, error) {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	list, err := db.New(r.db).GetFilesByBucketID(ctx, bucketID)
	if err != nil {
		return nil, err
	}
	out := make([]*db.File, 0, len(list))
	for i := range list {
		f := list[i]
		out = append(out, &f)
	}
	return out, nil
}

func (r *sqliteRepository) CreateFile(ctx context.Context, file *db.File) error {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	_, err := db.New(r.db).CreateFile(ctx, db.CreateFileParams{
		StringID:     file.StringID,
		BucketID:     file.BucketID,
		OriginalName: file.OriginalName,
		OwnerID:      file.OwnerID,
		Size:         file.Size,
		ContentType:  file.ContentType,
		S3Key:        file.S3Key,
		CreatedAt:    file.CreatedAt,
	})
	return err
}

func (r *sqliteRepository) UpdateFile(ctx context.Context, file *db.File) error {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	return db.New(r.db).UpdateFile(ctx, db.UpdateFileParams{
		OriginalName: file.OriginalName,
		OwnerID:      file.OwnerID,
		StringID:     file.StringID,
	})
}

func (r *sqliteRepository) DeleteFile(ctx context.Context, id int64) error {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	q := db.New(r.db)
	file, err := q.GetFileByID(ctx, id)
	if err != nil {
		return err
	}
	return q.DeleteFile(ctx, file.StringID)
}

func (r *sqliteRepository) ListFiles(ctx context.Context, limit int, offset int) ([]*db.File, error) {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	list, err := db.New(r.db).ListFiles(ctx, db.ListFilesParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, err
	}
	out := make([]*db.File, 0, len(list))
	for i := range list {
		f := list[i]
		out = append(out, &f)
	}
	return out, nil
}

// Bucket admin operations
func (r *sqliteRepository) AddBucketAdmin(ctx context.Context, bucketAdmin *db.BucketAdmin) error {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	return db.New(r.db).AddBucketAdmin(ctx, db.AddBucketAdminParams{
		UserID:    bucketAdmin.UserID,
		BucketID:  bucketAdmin.BucketID,
		CreatedAt: bucketAdmin.CreatedAt,
	})
}

func (r *sqliteRepository) RemoveBucketAdmin(ctx context.Context, userID string, bucketID string) error {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	return db.New(r.db).RemoveBucketAdmin(ctx, db.RemoveBucketAdminParams{
		UserID:   userID,
		BucketID: bucketID,
	})
}

func (r *sqliteRepository) GetBucketAdminsByBucketID(ctx context.Context, bucketID string) ([]*db.BucketAdmin, error) {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	list, err := db.New(r.db).GetBucketAdminsByBucketID(ctx, bucketID)
	if err != nil {
		return nil, err
	}
	out := make([]*db.BucketAdmin, 0, len(list))
	for i := range list {
		a := list[i]
		out = append(out, &a)
	}
	return out, nil
}

func (r *sqliteRepository) GetBucketsByAdminUserID(ctx context.Context, userID string) ([]*db.Bucket, error) {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	list, err := db.New(r.db).GetBucketsByAdminUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]*db.Bucket, 0, len(list))
	for i := range list {
		b := list[i]
		out = append(out, &b)
	}
	return out, nil
}

func (r *sqliteRepository) IsBucketAdmin(ctx context.Context, userID string, bucketID string) (bool, error) {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()
	v, err := db.New(r.db).IsBucketAdmin(ctx, db.IsBucketAdminParams{
		UserID:   userID,
		BucketID: bucketID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return v == 1, nil
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
	return context.WithTimeout(context.Background(), internalpkg.DEFAULT_REPOSITORY_QUERY_TIMEOUT)
}
