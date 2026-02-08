package pkg

import "time"

type Lifecycle struct {
	ID         int       `db:"id"`
	BucketSlug string    `db:"bucket_slug"`
	ExpiresAt  time.Time `db:"expires_at"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type PurgeExpiredBucketsResult struct {
	BucketSlug   string
	FilesDeleted int64
	Success      bool
}
