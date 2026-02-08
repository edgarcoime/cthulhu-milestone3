package handlers

import (
	"log/slog"
	"mime/multipart"
	"strings"
	"time"

	"github.com/cthulhu-platform/gateway/internal/connections"
	"github.com/cthulhu-platform/gateway/internal/middleware"
	"github.com/cthulhu-platform/gateway/internal/models"
	gatewaypkg "github.com/cthulhu-platform/gateway/internal/pkg"
	fmpb "github.com/cthulhu-platform/proto/pkg/filemanager"
	"github.com/gofiber/fiber/v2"
)

func parsePrepareUploadFromForm(c *fiber.Ctx) (*models.PrepareUploadRequest, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}
	defer form.RemoveAll()

	fileHeaders := form.File["files"]
	if len(fileHeaders) == 0 {
		return nil, nil
	}

	files := make([]models.PrepareUploadFile, 0, len(fileHeaders))
	for _, fh := range fileHeaders {
		ct := contentTypeFromHeader(fh)
		if ct == "" {
			ct = "application/octet-stream"
		}
		files = append(files, models.PrepareUploadFile{
			OriginalName: fh.Filename,
			Size:         fh.Size,
			ContentType:  ct,
		})
	}

	password := c.Get("X-Bucket-Password")
	if vs := form.Value["password"]; len(vs) > 0 && vs[0] != "" {
		password = vs[0]
	}

	return &models.PrepareUploadRequest{Files: files, Password: password}, nil
}

func contentTypeFromHeader(fh *multipart.FileHeader) string {
	if fh.Header == nil {
		return ""
	}
	ct := fh.Header.Get("Content-Type")
	return ct
}

func FileUploadPrepare(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		contentType := c.Get("Content-Type")
		var req models.PrepareUploadRequest

		if strings.HasPrefix(contentType, "application/json") {
			if err := c.BodyParser(&req); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
			}
		} else if strings.HasPrefix(contentType, "multipart/form-data") {
			parsed, err := parsePrepareUploadFromForm(c)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid form data"})
			}
			if parsed == nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "files is required and must be non-empty"})
			}
			req = *parsed
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Content-Type must be application/json or multipart/form-data"})
		}

		if len(req.Files) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "files is required and must be non-empty"})
		}

		pbFiles := make([]*fmpb.FileMeta, 0, len(req.Files))
		for i := range req.Files {
			f := &req.Files[i]
			if f.OriginalName == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "each file must have original_name"})
			}
			if f.Size < 0 {
				f.Size = 0
			}
			ct := f.ContentType
			if ct == "" {
				ct = "application/octet-stream"
			}
			pbFiles = append(pbFiles, &fmpb.FileMeta{
				OriginalName: f.OriginalName,
				Size:         f.Size,
				ContentType:  ct,
			})
		}

		var userID *string
		if u := middleware.GetUser(c); u != nil {
			userID = &u.ID
		}

		pbReq := &fmpb.PrepareUploadRequest{
			Files:  pbFiles,
			UserId: userID,
		}
		if req.Password != "" {
			pbReq.Password = &req.Password
		}

		res, err := conns.Filemanager.PrepareUpload(c.Context(), pbReq)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if res != nil && res.Error != "" {
			out := models.PrepareUploadResponse{Error: res.Error, StorageID: res.StorageId}
			return c.Status(fiber.StatusBadRequest).JSON(out)
		}

		slots := make([]models.UploadSlot, 0, len(res.Slots))
		for _, s := range res.Slots {
			slots = append(slots, models.UploadSlot{
				StringID:        s.StringId,
				PresignedPutURL: s.PresignedPutUrl,
				S3Key:           s.S3Key,
			})
		}
		return c.Status(fiber.StatusOK).JSON(models.PrepareUploadResponse{
			StorageID: res.StorageId,
			Slots:     slots,
		})
	}
}

func FileUploadConfirm(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.ConfirmUploadRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}
		if req.StorageID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "storage_id is required"})
		}
		if len(req.Files) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "files is required and must be non-empty"})
		}

		pbFiles := make([]*fmpb.FileMetaWithStringId, 0, len(req.Files))
		for i := range req.Files {
			f := &req.Files[i]
			if f.StringID == "" || f.OriginalName == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "each file must have string_id and original_name"})
			}
			pbFiles = append(pbFiles, &fmpb.FileMetaWithStringId{
				StringId:     f.StringID,
				OriginalName: f.OriginalName,
				Size:         f.Size,
				ContentType:  f.ContentType,
			})
		}

		pbReq := &fmpb.ConfirmUploadRequest{
			StorageId: req.StorageID,
			Files:     pbFiles,
		}

		res, err := conns.Filemanager.ConfirmUpload(c.Context(), pbReq)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if res != nil && res.Error != "" {
			out := models.ConfirmUploadResponse{Success: false, Error: res.Error, StorageID: res.StorageId}
			return c.Status(fiber.StatusBadRequest).JSON(out)
		}

		// Ensure bucket has a lifecycle (create or update). Best-effort: log and continue on failure.
		var expiresAt time.Time
		if middleware.GetUser(c) != nil {
			expiresAt = time.Now().UTC().Add(gatewaypkg.LifecycleTTLAuthorized)
		} else {
			expiresAt = time.Now().UTC().Add(gatewaypkg.LifecycleTTLAnonymous)
		}
		if _, err := conns.Lifecycle.PostLifecycle(c.Context(), res.StorageId, expiresAt); err != nil {
			slog.Warn("failed to set bucket lifecycle", "storage_id", res.StorageId, "error", err)
		}

		files := make([]models.FileInfoResult, 0, len(res.Files))
		for _, f := range res.Files {
			files = append(files, models.FileInfoResult{
				OriginalName: f.OriginalName,
				StringID:     f.StringId,
				Key:          f.Key,
				Size:         f.Size,
				ContentType:  f.ContentType,
			})
		}
		return c.Status(fiber.StatusOK).JSON(models.ConfirmUploadResponse{
			Success:   res.Success,
			StorageID: res.StorageId,
			Files:     files,
			TotalSize: res.TotalSize,
		})
	}
}

func FileUpload(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: complete RPC logic
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "uploading file",
		})
	}
}

func FileAuthenticate(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		bucketID := strings.TrimSpace(c.Params("id"))
		if bucketID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bucket id is required"})
		}
		var body struct {
			Password string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}
		if body.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "password is required"})
		}
		pbReq := &fmpb.AuthenticateBucketRequest{BucketId: bucketID, Password: body.Password}
		if u := middleware.GetUser(c); u != nil {
			pbReq.UserId = &u.ID
		}
		res, err := conns.Filemanager.AuthenticateBucket(c.Context(), pbReq)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if res != nil && res.Error != "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": res.Error})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"access_token": res.AccessToken,
			"expires_in":   res.ExpiresIn,
		})
	}
}

func FileBucketGet(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storageID := strings.TrimSpace(c.Params("id"))
		if storageID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "storage id is required"})
		}
		res, err := conns.Filemanager.RetrieveFileBucket(c.Context(), &fmpb.RetrieveFileBucketRequest{StorageId: storageID})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if res != nil && res.Error != "" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": res.Error})
		}
		files := make([]fiber.Map, 0, len(res.Files))
		for _, f := range res.Files {
			files = append(files, fiber.Map{
				"original_name": f.OriginalName,
				"string_id":     f.StringId,
				"key":           f.Key,
				"size":          f.Size,
				"content_type":  f.ContentType,
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"storage_id": res.StorageId,
			"files":      files,
			"total_size": res.TotalSize,
		})
	}
}

func FileAdmins(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		bucketID := strings.TrimSpace(c.Params("id"))
		if bucketID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bucket id is required"})
		}
		res, err := conns.Filemanager.GetBucketAdmins(c.Context(), &fmpb.GetBucketAdminsRequest{BucketId: bucketID})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if res != nil && res.Error != "" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": res.Error})
		}
		admins := make([]fiber.Map, 0, len(res.Admins))
		for _, a := range res.Admins {
			admins = append(admins, fiber.Map{
				"user_id":    a.UserId,
				"email":      a.Email,
				"username":   a.Username,
				"avatar_url": a.AvatarUrl,
				"is_owner":   a.IsOwner,
				"created_at": a.CreatedAt,
			})
		}
		out := fiber.Map{"bucket_id": res.BucketId, "admins": admins}
		if res.Owner != nil {
			out["owner"] = fiber.Map{
				"user_id":    res.Owner.UserId,
				"email":      res.Owner.Email,
				"username":   res.Owner.Username,
				"avatar_url": res.Owner.AvatarUrl,
				"is_owner":   res.Owner.IsOwner,
				"created_at": res.Owner.CreatedAt,
			}
		}
		return c.Status(fiber.StatusOK).JSON(out)
	}
}

func FileBucketProtected(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		bucketID := strings.TrimSpace(c.Params("id"))
		if bucketID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bucket id is required"})
		}
		res, err := conns.Filemanager.IsBucketProtected(c.Context(), &fmpb.IsBucketProtectedRequest{BucketId: bucketID})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if res != nil && res.Error != "" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": res.Error})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"protected": res.Protected,
			"bucket_id": bucketID,
		})
	}
}

func FileDownload(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storageID := strings.TrimSpace(c.Params("id"))
		stringID := strings.TrimSpace(c.Params("filename"))
		if storageID == "" || stringID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "storage id and file string_id are required"})
		}

		pbReq := &fmpb.PrepareDownloadRequest{
			StorageId: storageID,
			StringId:  stringID,
		}
		if token := c.Get("X-Bucket-Token"); token != "" {
			pbReq.BucketAccessToken = &token
		}

		res, err := conns.Filemanager.PrepareDownload(c.Context(), pbReq)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if res != nil && res.Error != "" {
			if res.Error == "bucket is protected; bucket_access_token is required" || res.Error == "invalid or expired bucket token" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": res.Error})
			}
			if res.Error == "file not found" {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": res.Error})
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": res.Error})
		}

		return c.Redirect(res.PresignedGetUrl, fiber.StatusFound)
	}
}
