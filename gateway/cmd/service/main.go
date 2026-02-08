package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/cthulhu-platform/gateway/internal/connections"
	internalpkg "github.com/cthulhu-platform/gateway/internal/pkg"
	"github.com/cthulhu-platform/gateway/internal/server"
)

func main() {
	// Setup logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx := context.Background()

	// Setup Dependencies (5s timeout for connection initialization)
	initCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	connectionPool, err := connections.NewConnectionsContainer(initCtx, connections.ConnectionsConfig{
		LifecycleURL:   internalpkg.LIFECYCLE_GRPC_URL,
		AuthURL:        internalpkg.AUTH_GRPC_URL,
		FilemanagerURL: internalpkg.FILEMANAGER_GRPC_URL,
	})
	if err != nil {
		slog.Error("Failed to create connections container", "error", err)
		os.Exit(1)
	}
	defer connectionPool.Close()

	serverCfg := server.FiberServerConfig{
		Host:   internalpkg.APP_HOST,
		Port:   internalpkg.APP_PORT,
		Logger: logger,
	}

	// Start HTTP server
	s := server.NewFiberServer(serverCfg, connectionPool)
	s.Start()
}
