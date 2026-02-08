package connections

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	filemanager "github.com/cthulhu-platform/filemanager/pkg/client"
)

type ConnectionsContainer struct {
	Filemanager *filemanager.Client
}

type ConnectionsConfig struct {
	FilemanagerURL string
}

func NewConnectionsContainer(ctx context.Context, cfg ConnectionsConfig) (*ConnectionsContainer, error) {
	// Create timeout Context
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Connect to Filemanager
	filemanagerClient, err := filemanager.NewClient(ctx, cfg.FilemanagerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create filemanager client: %v", err)
	}

	slog.Info("Filemanager client created", "url", cfg.FilemanagerURL)
	return &ConnectionsContainer{
		Filemanager: filemanagerClient,
	}, nil
}

func (c *ConnectionsContainer) Close() {
	c.Filemanager.Close()
}
