package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/cthulhu-platform/filemanager/internal/configs"
	"github.com/cthulhu-platform/filemanager/internal/connections"
	"github.com/cthulhu-platform/filemanager/internal/pkg"
	"github.com/cthulhu-platform/filemanager/internal/repository"
	"github.com/cthulhu-platform/filemanager/internal/server"
	"github.com/cthulhu-platform/filemanager/internal/service"
	"github.com/cthulhu-platform/filemanager/internal/storage"
)

func main() {
	fmt.Println("Filemanager service")

	ctx := context.Background()

	_, err := configs.Load("filemanager")
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Create dependencies

	// Create local Repository/db
	repo, err := repository.NewSQLiteRepository(ctx)
	if err != nil {
		slog.Error("Failed to create repository", "error", err)
		os.Exit(1)
	}
	defer repo.Close()

	// Create connection to storage capable service (AWS S3 compatible)
	presignedEndpoint := pkg.S3_PRESIGNED_ENDPOINT
	if presignedEndpoint == "" {
		presignedEndpoint = pkg.S3_ENDPOINT
	}
	storage, err := storage.NewAWSStorage(ctx, storage.AWSStorageConfig{
		AccessKeyID:        pkg.S3_ACCESS_KEY_ID,
		SecretAccessKey:    pkg.S3_SECRET_ACCESS_KEY,
		Endpoint:           pkg.S3_ENDPOINT,
		PresignedEndpoint:  presignedEndpoint,
		Region:             pkg.S3_REGION,
		BucketName:         pkg.S3_BUCKET_NAME,
		ForcePathStyle:     pkg.S3_FORCE_PATH_STYLE == "true",
	})
	if err != nil {
		slog.Error("Failed to connect to storage", "error", err)
		os.Exit(1)
	}
	defer storage.Close()

	// Create connections to other microservices
	connectionPool, err := connections.NewConnectionsContainer(ctx, connections.ConnectionsConfig{
		AuthURL: pkg.AUTH_GRPC_URL,
	})
	if err != nil {
		slog.Error("Failed to create connections container", "error", err)
		os.Exit(1)
	}
	defer connectionPool.Close()

	// Create Service (storage implements storage.Storage for PresignPut)
	svc := service.NewFilemanagerService(repo, storage, connectionPool)

	serverCfg := server.ServerConfig{
		Host: pkg.APP_HOST,
		Port: pkg.APP_PORT,
	}
	if err := server.ListenGRPC(ctx, serverCfg, svc); err != nil {
		slog.Error("Failed to start gRPC server", "error", err)
		os.Exit(1)
	}
}
