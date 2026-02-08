package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/cthulhu-platform/common/pkg/strings"
	"github.com/cthulhu-platform/lifecycle/internal/service"
	"github.com/cthulhu-platform/lifecycle/pkg"
	"github.com/cthulhu-platform/lifecycle/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type grpcServer struct {
	pb.UnimplementedLifecycleServiceServer
	service service.Service
}

type ServerConfig struct {
	Host string
	Port string
}

func ListenGRPC(ctx context.Context, cfg ServerConfig, svc service.Service) error {
	lis, err := net.Listen("tcp", cfg.Host+":"+cfg.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterLifecycleServiceServer(server, &grpcServer{service: svc})

	reflection.Register(server)
	slog.Info("Lifecycle service gRPC server listening", "host", cfg.Host, "port", cfg.Port)
	return server.Serve(lis)
}

func (s *grpcServer) PostLifecycle(ctx context.Context, req *pb.PostLifecycleRequest) (*pb.PostLifecycleResponse, error) {
	res, err := s.service.PostLifecycle(ctx, req.GetBucketSlug(), req.GetExpiresAt().AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not create lifecycle: %v", err)
	}
	slog.Info("Lifecycle created for bucket", "bucket_slug", strings.TruncateString(req.GetBucketSlug(), 4))
	return &pb.PostLifecycleResponse{Lifecycle: lifecycleToPB(res)}, nil
}

func (s *grpcServer) GetLifecycle(ctx context.Context, req *pb.GetLifecycleRequest) (*pb.GetLifecycleResponse, error) {
	res, err := s.service.GetLifecycle(ctx, req.GetBucketSlug())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not get retrieve lifecycle for bucket %s: %v", req.GetBucketSlug(), err)
	}
	slog.Info("Lifecycle retrieved for bucket", "bucket_slug", strings.TruncateString(req.GetBucketSlug(), 4))
	return &pb.GetLifecycleResponse{Lifecycle: lifecycleToPB(res)}, nil
}

func (s *grpcServer) DeleteLifecycle(ctx context.Context, req *pb.DeleteLifecycleRequest) (*pb.DeleteLifecycleResponse, error) {
	err := s.service.DeleteLifecycle(ctx, req.GetBucketSlug())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not delete lifecycle for bucket %s: %v", req.GetBucketSlug(), err)
	}
	slog.Info("Lifecycle deleted for bucket", "bucket_slug", strings.TruncateString(req.GetBucketSlug(), 4))
	return &pb.DeleteLifecycleResponse{Success: true}, nil
}

func lifecycleToPB(l *pkg.Lifecycle) *pb.Lifecycle {
	if l == nil {
		return nil
	}
	return &pb.Lifecycle{
		Id:         int32(l.ID),
		BucketSlug: l.BucketSlug,
		ExpiresAt:  timestamppb.New(l.ExpiresAt),
		CreatedAt:  timestamppb.New(l.CreatedAt),
		UpdatedAt:  timestamppb.New(l.UpdatedAt),
	}
}
