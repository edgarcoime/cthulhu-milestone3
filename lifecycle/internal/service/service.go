package service

import (
	"context"
	"time"

	"github.com/cthulhu-platform/lifecycle/internal/connections"
	"github.com/cthulhu-platform/lifecycle/internal/repository"
	"github.com/cthulhu-platform/lifecycle/pkg"
	pb "github.com/cthulhu-platform/proto/pkg/filemanager"
)

type Service interface {
	PostLifecycle(ctx context.Context, bucketSlug string, expiresAt time.Time) (*pkg.Lifecycle, error)
	GetLifecycle(ctx context.Context, bucketSlug string) (*pkg.Lifecycle, error)
	DeleteLifecycle(ctx context.Context, bucketSlug string) error
	PurgeExpiredBuckets(ctx context.Context) ([]pkg.PurgeExpiredBucketsResult, error)
}

type lifecycleService struct {
	repo  repository.Repository
	conns *connections.ConnectionsContainer
}

func NewLifecycleService(repo repository.Repository, conns *connections.ConnectionsContainer) Service {
	return &lifecycleService{repo: repo, conns: conns}
}

func (s *lifecycleService) PostLifecycle(ctx context.Context, bucketSlug string, expiresAt time.Time) (*pkg.Lifecycle, error) {
	return s.repo.PutLifecycle(ctx, pkg.Lifecycle{BucketSlug: bucketSlug, ExpiresAt: expiresAt})
}

func (s *lifecycleService) GetLifecycle(ctx context.Context, bucketSlug string) (*pkg.Lifecycle, error) {
	return s.repo.GetLifecycle(ctx, bucketSlug)
}

func (s *lifecycleService) DeleteLifecycle(ctx context.Context, bucketSlug string) error {
	return s.repo.DeleteLifecycle(ctx, bucketSlug)
}

func (s *lifecycleService) PurgeExpiredBuckets(ctx context.Context) ([]pkg.PurgeExpiredBucketsResult, error) {
	// Check to see if shown first
	now := time.Now().UTC()
	expired, err := s.repo.ListExpiredLifecycles(ctx, now)
	if err != nil {
		return nil, err
	}
	results := make([]pkg.PurgeExpiredBucketsResult, 0, len(expired))

	for _, l := range expired {
		result := pkg.PurgeExpiredBucketsResult{BucketSlug: l.BucketSlug}
		resp, err := s.conns.Filemanager.DeleteBucket(ctx, &pb.DeleteBucketRequest{BucketId: l.BucketSlug})
		if err != nil || (resp != nil && resp.Error != "") {
			result.Success = false
			results = append(results, result)
			continue
		}
		result.FilesDeleted = resp.FilesDeleted
		result.Success = resp.Success
		_ = s.repo.DeleteLifecycle(ctx, l.BucketSlug)
		results = append(results, result)
	}

	return results, nil
}
