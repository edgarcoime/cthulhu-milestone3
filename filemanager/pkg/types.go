package pkg

import (
	"io"

	"github.com/golang-jwt/jwt/v5"
)

type AdminInfo struct {
	UserID    string  `json:"user_id"`
	Email     string  `json:"email"`
	Username  *string `json:"username,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	IsOwner   bool    `json:"is_owner"`
	CreatedAt int64   `json:"created_at"`
}

type BucketAdminsResponse struct {
	BucketID string      `json:"bucket_id"`
	Owner    *AdminInfo  `json:"owner"`
	Admins   []AdminInfo `json:"admins"`
}

// FileInfo represents a stored object.
type FileInfo struct {
	OriginalName string `json:"original_name"`
	StringID     string `json:"string_id"`
	Key          string `json:"key"`
	Size         int64  `json:"size"`
	ContentType  string `json:"content_type"`
}

// UploadResult is returned after an upload transaction.
type UploadResult struct {
	TransactionID string     `json:"transaction_id"`
	Success       bool       `json:"success"`
	Error         string     `json:"error,omitempty"`
	StorageID     string     `json:"storage_id,omitempty"`
	Files         []FileInfo `json:"files,omitempty"`
	TotalSize     int64      `json:"total_size,omitempty"`
}

// BucketMetadata contains objects under a storage ID.
type BucketMetadata struct {
	StorageID string     `json:"storage_id"`
	Files     []FileInfo `json:"files"`
	TotalSize int64      `json:"total_size"`
}

// DownloadResult wraps object body and metadata for streaming.
type DownloadResult struct {
	Body           io.ReadCloser `json:"body"`
	ContentType    string        `json:"content_type"`
	ContentLength  int64         `json:"content_length"`
	DownloadedFile string        `json:"downloaded_file"`
}

// UploadObject is a single file to upload.
type UploadObject struct {
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	Body        io.Reader `json:"body"`
}

// TOKEN RELATED TYPES

// BucketAccessClaims represents JWT claims for bucket access tokens
type BucketAccessClaims struct {
	BucketID    string   `json:"bucket_id"`
	Privileges  []string `json:"privileges"` // ["read", "write", etc.]
	UserID      *string  `json:"user_id,omitempty"`
	AuthTokenID *string  `json:"auth_token_id,omitempty"` // JTI from auth token
	jwt.RegisteredClaims
}

// BucketAccessTokenResponse represents the API response for bucket authentication
type BucketAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"` // seconds
}
