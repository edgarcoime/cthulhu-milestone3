package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/cthulhu-platform/auth/internal/pkg"
	"github.com/cthulhu-platform/auth/internal/repository"
	"github.com/cthulhu-platform/auth/internal/server"
	"github.com/cthulhu-platform/auth/internal/service"
)

func main() {
	fmt.Println("Auth Service")

	ctx := context.Background()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	repo, err := repository.NewSQLiteRepository(ctx)
	if err != nil {
		logger.Error("Failed to create repository", "error", err)
		os.Exit(1)
	}
	defer repo.Close()

	svc := service.NewAuthService(repo)

	serverCfg := server.ServerConfig{
		Host: pkg.APP_HOST,
		Port: pkg.APP_PORT,
	}

	err = server.ListenGRPC(ctx, serverCfg, svc)
	if err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
