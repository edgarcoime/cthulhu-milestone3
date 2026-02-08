package models

// PrepareUpload (request)

type PrepareUploadFile struct {
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	ContentType  string `json:"content_type"`
}

type PrepareUploadRequest struct {
	Files    []PrepareUploadFile `json:"files"`
	Password string              `json:"password,omitempty"`
}

// ConfirmUpload (request)

type ConfirmUploadFile struct {
	StringID     string `json:"string_id"`
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	ContentType  string `json:"content_type"`
}

type ConfirmUploadRequest struct {
	StorageID string              `json:"storage_id"`
	Files     []ConfirmUploadFile `json:"files"`
}
