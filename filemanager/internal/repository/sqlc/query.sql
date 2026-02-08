-- Buckets

-- name: GetBucketByID :one
SELECT * FROM buckets WHERE id = ? LIMIT 1;

-- name: CreateBucket :exec
INSERT INTO buckets (id, password_hash, created_at, updated_at)
VALUES (?, ?, ?, ?);

-- name: UpdateBucket :exec
UPDATE buckets SET password_hash = ?, updated_at = ? WHERE id = ?;

-- name: DeleteBucket :exec
DELETE FROM buckets WHERE id = ?;

-- name: ListBuckets :many
SELECT * FROM buckets ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- Files

-- name: GetFileByStringID :one
SELECT * FROM files WHERE string_id = ? LIMIT 1;

-- name: GetFileByBucketIDAndStringID :one
SELECT * FROM files WHERE bucket_id = ? AND string_id = ? LIMIT 1;

-- name: GetFileByID :one
SELECT * FROM files WHERE id = ? LIMIT 1;

-- name: GetFilesByBucketID :many
SELECT * FROM files WHERE bucket_id = ? ORDER BY created_at ASC;

-- name: GetFilesByOwnerID :many
SELECT * FROM files WHERE owner_id = ? ORDER BY created_at DESC;

-- name: CreateFile :one
INSERT INTO files (string_id, bucket_id, original_name, owner_id, size, content_type, s3_key, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateFile :exec
UPDATE files SET original_name = ?, owner_id = ? WHERE string_id = ?;

-- name: DeleteFile :exec
DELETE FROM files WHERE string_id = ?;

-- name: CountFilesByBucketID :one
SELECT COUNT(*) FROM files WHERE bucket_id = ?;

-- name: ListFiles :many
SELECT * FROM files ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- Bucket admins

-- name: AddBucketAdmin :exec
INSERT INTO bucket_admins (user_id, bucket_id, created_at)
VALUES (?, ?, ?);

-- name: RemoveBucketAdmin :exec
DELETE FROM bucket_admins WHERE user_id = ? AND bucket_id = ?;

-- name: GetBucketAdminsByBucketID :many
SELECT * FROM bucket_admins WHERE bucket_id = ? ORDER BY created_at ASC;

-- name: GetBucketsByAdminUserID :many
SELECT b.* FROM buckets b
INNER JOIN bucket_admins ba ON b.id = ba.bucket_id
WHERE ba.user_id = ?;

-- name: IsBucketAdmin :one
SELECT 1 FROM bucket_admins WHERE user_id = ? AND bucket_id = ? LIMIT 1;
