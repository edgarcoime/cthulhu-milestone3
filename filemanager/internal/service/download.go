// PrepareDownload returns a presigned GET URL for direct S3 download.
// For protected buckets, bucket_access_token (from AuthenticateBucket) is required.

package service

import (
	"context"
	"errors"

	pb "github.com/cthulhu-platform/proto/pkg/filemanager"
)

func (s *filemanagerService) PrepareDownload(ctx context.Context, req *pb.PrepareDownloadRequest) (*pb.PrepareDownloadResponse, error) {
	res := &pb.PrepareDownloadResponse{}
	if req == nil || req.StorageId == "" || req.StringId == "" {
		res.Error = "storage_id and string_id are required"
		return res, errors.New(res.Error)
	}

	storageID := req.StorageId
	stringID := req.StringId

	file, err := s.repo.GetFileByBucketIDAndStringID(ctx, storageID, stringID)
	if err != nil {
		res.Error = "file not found"
		return res, err
	}

	bucket, err := s.repo.GetBucketByID(ctx, storageID)
	if err != nil {
		res.Error = err.Error()
		return res, err
	}

	if bucket.PasswordHash.Valid {
		token := ""
		if req.BucketAccessToken != nil {
			token = *req.BucketAccessToken
		}
		if token == "" {
			res.Error = "bucket is protected; bucket_access_token is required"
			return res, errors.New(res.Error)
		}
		claims, err := ValidateBucketAccessToken(token)
		if err != nil {
			res.Error = "invalid or expired bucket token"
			return res, err
		}
		if claims.BucketID != storageID {
			res.Error = "bucket token does not match bucket"
			return res, errors.New(res.Error)
		}
	}

	url, err := s.storage.PresignGet(ctx, file.S3Key)
	if err != nil {
		res.Error = err.Error()
		return res, err
	}

	res.PresignedGetUrl = url
	res.OriginalName = file.OriginalName
	res.ContentType = file.ContentType
	res.Size = file.Size
	return res, nil
}
