package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/cthulhu-platform/lifecycle/internal/connections"
	"github.com/cthulhu-platform/lifecycle/internal/daemon"
	internalpkg "github.com/cthulhu-platform/lifecycle/internal/pkg"
	"github.com/cthulhu-platform/lifecycle/internal/repository"
	"github.com/cthulhu-platform/lifecycle/internal/server"
	"github.com/cthulhu-platform/lifecycle/internal/service"
)

func main() {
	fmt.Println("Lifecycle service")

	ctx := context.Background()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	repo, err := repository.NewSQLiteRepository(ctx)
	if err != nil {
		slog.Error("Failed to create repository", "error", err)
		os.Exit(1)
	}
	defer repo.Close()

	// Create connections to other microservices
	conns, err := connections.NewConnectionsContainer(ctx, connections.ConnectionsConfig{
		FilemanagerURL: internalpkg.FILEMANAGER_GRPC_URL,
	})
	if err != nil {
		slog.Error("Failed to create connections container", "error", err)
		os.Exit(1)
	}
	defer conns.Close()

	svc := service.NewLifecycleService(repo, conns)

	serverCfg := server.ServerConfig{
		Host: internalpkg.APP_HOST,
		Port: internalpkg.APP_PORT,
	}

	cleanupDaemon := daemon.NewCleanupDaemon(repo, svc, internalpkg.DEFAULT_CLEANUP_INTERVAL)
	go cleanupDaemon.Run(ctx)

	err = server.ListenGRPC(ctx, serverCfg, svc)
	if err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
