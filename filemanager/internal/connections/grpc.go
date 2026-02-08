package connections

import (
	"context"
	"fmt"
	"time"

	auth "github.com/cthulhu-platform/auth/pkg/client"
)

type ConnectionsContainer struct {
	Auth *auth.Client
}

type ConnectionsConfig struct {
	AuthURL string
}

func NewConnectionsContainer(ctx context.Context, cfg ConnectionsConfig) (*ConnectionsContainer, error) {
	// Create timeout Context
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Connect to Auth
	authClient, err := auth.NewClient(ctx, cfg.AuthURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth client: %v", err)
	}

	return &ConnectionsContainer{
		Auth: authClient,
	}, nil
}

func (c *ConnectionsContainer) Close() {
	c.Auth.Close()
}
