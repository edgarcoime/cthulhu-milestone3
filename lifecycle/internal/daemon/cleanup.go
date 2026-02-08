package daemon

import (
	"context"
	"log/slog"
	"time"

	"github.com/cthulhu-platform/lifecycle/internal/repository"
	"github.com/cthulhu-platform/lifecycle/internal/service"
)

// Cleanup daemon that runs every interval and deletes expired lifecycles

type CleanupDaemon struct {
	repo     repository.Repository
	interval time.Duration
	service  service.Service
}

func NewCleanupDaemon(repo repository.Repository, service service.Service, interval time.Duration) *CleanupDaemon {
	return &CleanupDaemon{repo: repo, service: service, interval: interval}
}

func (d *CleanupDaemon) purge(ctx context.Context) {
	slog.Info("Bucket purge check started")
	results, err := d.service.PurgeExpiredBuckets(ctx)
	if err != nil {
		slog.Error("Purge expired buckets failed", "error", err)
		return
	}
	var successCount int
	for _, r := range results {
		if r.Success {
			successCount++
		}
	}
	if len(results) > 0 {
		slog.Info("Purged expired buckets", "total", len(results), "success", successCount)
	} else {
		slog.Info("No expired buckets to purge")
	}
	slog.Info("Bucket purge check completed")
}

func (d *CleanupDaemon) Run(ctx context.Context) error {
	slog.Info("Starting cleanup daemon", "interval", d.interval.String())
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()

	// Only add this if want to purge immediately on startup
	d.purge(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			d.purge(ctx)
		}
	}
}
