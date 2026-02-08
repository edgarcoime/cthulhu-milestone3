package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/cthulhu-platform/auth/internal/service"
	"github.com/cthulhu-platform/auth/pkg"
	"github.com/cthulhu-platform/common/pkg/strings"
	pb "github.com/cthulhu-platform/proto/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type grpcServer struct {
	pb.UnimplementedAuthServiceServer
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

	srv := grpc.NewServer()
	pb.RegisterAuthServiceServer(srv, &grpcServer{service: svc})

	reflection.Register(srv)
	slog.Info("Auth service gRPC server listening", "host", cfg.Host, "port", cfg.Port)
	return srv.Serve(lis)
}

func (s *grpcServer) InitiateOAuth(ctx context.Context, req *pb.InitiateOAuthRequest) (*pb.InitiateOAuthResponse, error) {
	redirectURL, err := s.service.InitiateOAuth(ctx, req.GetProvider())
	if err != nil {
		slog.Error("Failed to initiate OAuth", "error", err)
		return nil, status.Errorf(codes.Internal, "initiate OAuth: %v", err)
	}
	slog.Info("OAuth initiated for provider", "provider", req.GetProvider())
	return &pb.InitiateOAuthResponse{RedirectUrl: redirectURL}, nil
}

func (s *grpcServer) HandleOAuthCallback(ctx context.Context, req *pb.HandleOAuthCallbackRequest) (*pb.HandleOAuthCallbackResponse, error) {
	res, err := s.service.HandleOAuthCallback(ctx, req.GetProvider(), req.GetCode(), req.GetState())
	if err != nil {
		slog.Error("Failed to handle OAuth callback", "error", err)
		return nil, status.Errorf(codes.Internal, "handle OAuth callback: %v", err)
	}
	slog.Info("OAuth callback handled for provider", "provider", req.GetProvider())
	return &pb.HandleOAuthCallbackResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		User:         userInfoToPB(res.User),
	}, nil
}

func (s *grpcServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	user, err := s.service.ValidateToken(ctx, req.GetToken())
	if err != nil {
		slog.Error("Failed to validate token", "error", err)
		return nil, status.Errorf(codes.Internal, "validate token: %v", err)
	}
	slog.Info("Token validated for user", "user_id", strings.TruncateString(user.ID, 4))
	return &pb.ValidateTokenResponse{User: userInfoToPB(user)}, nil
}

func (s *grpcServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	res, err := s.service.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		slog.Error("Failed to refresh token", "error", err)
		return nil, status.Errorf(codes.Internal, "refresh token: %v", err)
	}
	slog.Info("Token refreshed", "refresh_token", strings.TruncateString(req.GetRefreshToken(), 4))
	return &pb.RefreshTokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (s *grpcServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	err := s.service.Logout(ctx, req.GetAccessToken())
	if err != nil {
		slog.Error("Failed to logout", "error", err)
		return nil, status.Errorf(codes.Internal, "logout: %v", err)
	}
	slog.Info("Token revoked", "access_token", strings.TruncateString(req.GetAccessToken(), 4))
	return &pb.LogoutResponse{Success: true}, nil
}

func userInfoToPB(u *pkg.UserInfo) *pb.UserInfo {
	if u == nil {
		return nil
	}
	return &pb.UserInfo{
		Id:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		AvatarUrl: u.AvatarUrl,
	}
}
