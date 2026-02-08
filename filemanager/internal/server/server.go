package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	internalpkg "github.com/cthulhu-platform/filemanager/internal/pkg"
	"github.com/cthulhu-platform/filemanager/internal/service"
	"github.com/cthulhu-platform/filemanager/pkg"

	// "github.com/cthulhu-platform/filemanager/pkg/pb"

	pb "github.com/cthulhu-platform/proto/pkg/filemanager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type grpcServer struct {
	pb.UnimplementedFilemanagerServiceServer
	svc service.Service
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
	pb.RegisterFilemanagerServiceServer(srv, &grpcServer{svc: svc})

	reflection.Register(srv)
	slog.Info("Filemanager gRPC server listening", "host", cfg.Host, "port", cfg.Port)
	return srv.Serve(lis)
}

func (s *grpcServer) PrepareUpload(ctx context.Context, req *pb.PrepareUploadRequest) (*pb.PrepareUploadResponse, error) {
	res, err := s.svc.PrepareUpload(ctx, req)
	if err != nil {
		// NOTE: Check is if the error is already set in the response
		// Such as something like "failed to generate unique storage id"
		if res != nil && res.Error != "" {
			return res, nil
		}
		return nil, status.Errorf(codes.Internal, "prepare upload: %v", err)
	}
	slog.Info("Prepare upload response", "storage_id", res.StorageId, "slots", len(res.Slots))
	return res, nil
}

func (s *grpcServer) ConfirmUpload(ctx context.Context, req *pb.ConfirmUploadRequest) (*pb.ConfirmUploadResponse, error) {
	res, err := s.svc.ConfirmUpload(ctx, req)
	if err != nil {
		if res != nil && res.Error != "" {
			return res, nil
		}
		return nil, status.Errorf(codes.Internal, "confirm upload: %v", err)
	}
	slog.Info("Confirm upload response", "storage_id", res.StorageId, "files", len(res.Files), "total_size", res.TotalSize)
	return res, nil
}

func (s *grpcServer) PrepareDownload(ctx context.Context, req *pb.PrepareDownloadRequest) (*pb.PrepareDownloadResponse, error) {
	res, err := s.svc.PrepareDownload(ctx, req)
	if err != nil {
		if res != nil && res.Error != "" {
			return res, nil
		}
		return nil, status.Errorf(codes.Internal, "prepare download: %v", err)
	}
	slog.Info("Prepare download response", "storage_id", req.StorageId, "original_name", res.OriginalName, "content_type", res.ContentType, "size", res.Size)
	return res, nil
}

func (s *grpcServer) RetrieveFileBucket(ctx context.Context, req *pb.RetrieveFileBucketRequest) (*pb.RetrieveFileBucketResponse, error) {
	meta, err := s.svc.RetrieveFileBucket(ctx, req.StorageId)
	if err != nil {
		return &pb.RetrieveFileBucketResponse{Error: err.Error()}, nil
	}
	out := &pb.RetrieveFileBucketResponse{
		StorageId: meta.StorageID,
		TotalSize: meta.TotalSize,
		Files:     make([]*pb.FileInfoResult, 0, len(meta.Files)),
	}
	for _, f := range meta.Files {
		out.Files = append(out.Files, &pb.FileInfoResult{
			OriginalName: f.OriginalName,
			StringId:     f.StringID,
			Key:          f.Key,
			Size:         f.Size,
			ContentType:  f.ContentType,
		})
	}
	slog.Info("Retrieve file bucket response", "storage_id", req.StorageId, "files", len(out.Files), "total_size", out.TotalSize)
	return out, nil
}

func adminInfoToPB(a pkg.AdminInfo) *pb.AdminInfo {
	return &pb.AdminInfo{
		UserId:    a.UserID,
		Email:     a.Email,
		Username:  a.Username,
		AvatarUrl: a.AvatarURL,
		IsOwner:   a.IsOwner,
		CreatedAt: a.CreatedAt,
	}
}

func (s *grpcServer) GetBucketAdmins(ctx context.Context, req *pb.GetBucketAdminsRequest) (*pb.GetBucketAdminsResponse, error) {
	admins, err := s.svc.GetBucketAdmins(ctx, req.BucketId)
	if err != nil {
		return &pb.GetBucketAdminsResponse{Error: err.Error()}, nil
	}
	out := &pb.GetBucketAdminsResponse{
		BucketId: admins.BucketID,
		Admins:   make([]*pb.AdminInfo, 0, len(admins.Admins)),
	}
	if admins.Owner != nil {
		out.Owner = adminInfoToPB(*admins.Owner)
	}
	for _, a := range admins.Admins {
		out.Admins = append(out.Admins, adminInfoToPB(a))
	}
	slog.Info("Get bucket admins response", "bucket_id", req.BucketId, "admins", len(out.Admins))
	return out, nil
}

func (s *grpcServer) IsBucketProtected(ctx context.Context, req *pb.IsBucketProtectedRequest) (*pb.IsBucketProtectedResponse, error) {
	protected, _, err := s.svc.IsBucketProtected(ctx, req.BucketId)
	if err != nil {
		return &pb.IsBucketProtectedResponse{Error: err.Error()}, nil
	}
	slog.Info("Is bucket protected response", "bucket_id", req.BucketId, "protected", protected)
	return &pb.IsBucketProtectedResponse{Protected: protected}, nil
}

func (s *grpcServer) AuthenticateBucket(ctx context.Context, req *pb.AuthenticateBucketRequest) (*pb.AuthenticateBucketResponse, error) {
	var userID, authTokenID *string
	if req.UserId != nil && *req.UserId != "" {
		userID = req.UserId
	}
	if req.AuthTokenId != nil && *req.AuthTokenId != "" {
		authTokenID = req.AuthTokenId
	}
	token, err := s.svc.AuthenticateBucket(ctx, req.BucketId, req.Password, userID, authTokenID)
	if err != nil {
		return &pb.AuthenticateBucketResponse{Error: err.Error()}, nil
	}
	expiresIn := int32(internalpkg.BUCKET_TOKEN_EXPIRATION.Seconds())
	slog.Info("Authenticate bucket response", "bucket_id", req.BucketId, "expires_in", expiresIn)
	return &pb.AuthenticateBucketResponse{AccessToken: token, ExpiresIn: expiresIn}, nil
}

func (s *grpcServer) DeleteBucket(ctx context.Context, req *pb.DeleteBucketRequest) (*pb.DeleteBucketResponse, error) {
	filesDeleted, err := s.svc.DeleteBucket(ctx, req.BucketId)
	if err != nil {
		return &pb.DeleteBucketResponse{Success: false, Error: err.Error()}, nil
	}
	slog.Info("Delete bucket response", "bucket_id", req.BucketId, "files_deleted", filesDeleted)
	return &pb.DeleteBucketResponse{Success: true, FilesDeleted: filesDeleted}, nil
}
