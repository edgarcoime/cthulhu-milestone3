package models

// PrepareUpload (response)

type UploadSlot struct {
	StringID        string `json:"string_id"`
	PresignedPutURL string `json:"presigned_put_url"`
	S3Key           string `json:"s3_key"`
}

type PrepareUploadResponse struct {
	StorageID string       `json:"storage_id,omitempty"`
	Slots     []UploadSlot `json:"slots,omitempty"`
	Error     string       `json:"error,omitempty"`
}

// ConfirmUpload (response)

type FileInfoResult struct {
	OriginalName string `json:"original_name"`
	StringID     string `json:"string_id"`
	Key          string `json:"key"`
	Size         int64  `json:"size"`
	ContentType  string `json:"content_type"`
}

type ConfirmUploadResponse struct {
	Success   bool             `json:"success"`
	StorageID string           `json:"storage_id,omitempty"`
	Files     []FileInfoResult `json:"files,omitempty"`
	TotalSize int64            `json:"total_size,omitempty"`
	Error     string           `json:"error,omitempty"`
}
