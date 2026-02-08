package connections

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	auth "github.com/cthulhu-platform/auth/pkg/client"
	filemanager "github.com/cthulhu-platform/filemanager/pkg/client"
	lifecycle "github.com/cthulhu-platform/lifecycle/pkg/client"
)

type ConnectionsContainer struct {
	Lifecycle   *lifecycle.Client
	Auth        *auth.Client
	Filemanager *filemanager.Client
}

type ConnectionsConfig struct {
	LifecycleURL   string
	AuthURL        string
	FilemanagerURL string
}

func NewConnectionsContainer(ctx context.Context, cfg ConnectionsConfig) (*ConnectionsContainer, error) {
	// Create a timeout context
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Connect to lifecycle service
	lifecycleClient, err := lifecycle.NewClient(ctx, cfg.LifecycleURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create lifecycle client: %v", err)
	}
	slog.Info("Lifecycle client created", "url", cfg.LifecycleURL)

	// Connect to auth service
	authClient, err := auth.NewClient(ctx, cfg.AuthURL)
	if err != nil {
		lifecycleClient.Close()
		return nil, fmt.Errorf("failed to create auth client: %v", err)
	}
	slog.Info("Auth client created", "url", cfg.AuthURL)

	// Connect to filemanager service
	filemanagerClient, err := filemanager.NewClient(ctx, cfg.FilemanagerURL)
	if err != nil {
		authClient.Close()
		lifecycleClient.Close()
		return nil, fmt.Errorf("failed to create filemanager client: %v", err)
	}
	slog.Info("Filemanager client created", "url", cfg.FilemanagerURL)

	return &ConnectionsContainer{
		Lifecycle:   lifecycleClient,
		Auth:        authClient,
		Filemanager: filemanagerClient,
	}, nil
}

func (c *ConnectionsContainer) Close() {
	c.Lifecycle.Close()
	c.Auth.Close()
	c.Filemanager.Close()
}
