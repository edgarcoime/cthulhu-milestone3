// Two-phase presigned URL upload: PrepareUpload returns presigned PUT URLs per file;
// the client uploads each file directly to S3; ConfirmUpload persists file metadata.

package service

import (
	"context"
	"database/sql"
	"errors"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/cthulhu-platform/filemanager/internal/repository/sqlc/db"
	pb "github.com/cthulhu-platform/proto/pkg/filemanager"
	"github.com/google/uuid"
)

func generateStorageID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 10
	var sb strings.Builder
	sb.Grow(length)
	for i := 0; i < length; i++ {
		sb.WriteByte(letters[rand.IntN(len(letters))])
	}
	return sb.String()
}

func (s *filemanagerService) PrepareUpload(ctx context.Context, req *pb.PrepareUploadRequest) (*pb.PrepareUploadResponse, error) {
	res := &pb.PrepareUploadResponse{}
	if req == nil || len(req.Files) == 0 {
		res.Error = "no files provided"
		return res, errors.New(res.Error)
	}

	storageID := generateStorageID()
	for i := 0; i < 3; i++ {
		existing, err := s.repo.GetBucketByID(ctx, storageID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			res.Error = err.Error()
			return res, err
		}
		if err != nil || existing == nil {
			// No row = ID is available; use it
			break
		}
		storageID = generateStorageID()
		if i == 2 {
			res.Error = "failed to generate unique storage id"
			return res, errors.New(res.Error)
		}
	}

	var passwordHash sql.NullString
	if req.Password != nil && *req.Password != "" {
		hash, err := HashBucketPassword(*req.Password)
		if err != nil {
			res.Error = err.Error()
			return res, err
		}
		passwordHash = sql.NullString{String: hash, Valid: true}
	}

	now := time.Now().Unix()
	bucket := &db.Bucket{
		ID:           storageID,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.repo.CreateBucket(ctx, bucket); err != nil {
		res.Error = err.Error()
		return res, err
	}

	if req.UserId != nil && *req.UserId != "" {
		admin := &db.BucketAdmin{
			UserID:    *req.UserId,
			BucketID:  storageID,
			CreatedAt: now,
		}
		_ = s.repo.AddBucketAdmin(ctx, admin)
	}

	slots := make([]*pb.FileUploadSlot, 0, len(req.Files))
	for _, f := range req.Files {
		stringID := uuid.New().String()
		s3Key := storageID + "/" + stringID
		size := f.Size
		if size < 0 {
			size = 0
		}
		contentType := f.ContentType
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		url, err := s.storage.PresignPut(ctx, s3Key, size, contentType)
		if err != nil {
			res.StorageId = storageID
			res.Error = err.Error()
			return res, err
		}
		slots = append(slots, &pb.FileUploadSlot{
			StringId:        stringID,
			PresignedPutUrl: url,
			S3Key:           s3Key,
		})
	}

	res.StorageId = storageID
	res.Slots = slots
	return res, nil
}

func (s *filemanagerService) ConfirmUpload(ctx context.Context, req *pb.ConfirmUploadRequest) (*pb.ConfirmUploadResponse, error) {
	res := &pb.ConfirmUploadResponse{Success: false}
	if req == nil || req.StorageId == "" || len(req.Files) == 0 {
		res.Error = "storage_id and at least one file required"
		return res, errors.New(res.Error)
	}

	_, err := s.repo.GetBucketByID(ctx, req.StorageId)
	if err != nil {
		res.Error = err.Error()
		return res, err
	}

	now := time.Now().Unix()
	var totalSize int64
	fileResults := make([]*pb.FileInfoResult, 0, len(req.Files))
	for _, f := range req.Files {
		s3Key := req.StorageId + "/" + f.StringId
		ownerID := sql.NullString{Valid: false}
		dbFile := &db.File{
			StringID:     f.StringId,
			BucketID:     req.StorageId,
			OriginalName: f.OriginalName,
			OwnerID:      ownerID,
			Size:         f.Size,
			ContentType:  f.ContentType,
			S3Key:        s3Key,
			CreatedAt:    now,
		}
		if err := s.repo.CreateFile(ctx, dbFile); err != nil {
			res.StorageId = req.StorageId
			res.Error = err.Error()
			return res, err
		}
		totalSize += f.Size
		fileResults = append(fileResults, &pb.FileInfoResult{
			OriginalName: f.OriginalName,
			StringId:     f.StringId,
			Key:          s3Key,
			Size:         f.Size,
			ContentType:  f.ContentType,
		})
	}

	res.Success = true
	res.StorageId = req.StorageId
	res.Files = fileResults
	res.TotalSize = totalSize
	return res, nil
}
