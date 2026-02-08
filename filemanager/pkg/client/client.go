package client

import (
	"context"
	"fmt"

	pb "github.com/cthulhu-platform/proto/pkg/filemanager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client is a gRPC client for the filemanager service.
type Client struct {
	conn    *grpc.ClientConn
	service pb.FilemanagerServiceClient
}

// NewClient creates a new filemanager gRPC client.
func NewClient(ctx context.Context, addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	return &Client{
		conn:    conn,
		service: pb.NewFilemanagerServiceClient(conn),
	}, nil
}

// Close closes the client connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// NOTE: This client is basically a wrapper for the service interface
// Should we remove PB implementation directly in the service and do conversion here?

// PrepareUpload requests presigned PUT URLs for uploading files.
func (c *Client) PrepareUpload(ctx context.Context, req *pb.PrepareUploadRequest) (*pb.PrepareUploadResponse, error) {
	return c.service.PrepareUpload(ctx, req)
}

// ConfirmUpload confirms that uploads are complete and persists file metadata.
func (c *Client) ConfirmUpload(ctx context.Context, req *pb.ConfirmUploadRequest) (*pb.ConfirmUploadResponse, error) {
	return c.service.ConfirmUpload(ctx, req)
}

// PrepareDownload returns a presigned GET URL for direct S3 download. For protected buckets, include bucket_access_token in the request.
func (c *Client) PrepareDownload(ctx context.Context, req *pb.PrepareDownloadRequest) (*pb.PrepareDownloadResponse, error) {
	return c.service.PrepareDownload(ctx, req)
}

// RetrieveFileBucket returns bucket metadata (files and total size) for the given storage ID.
func (c *Client) RetrieveFileBucket(ctx context.Context, req *pb.RetrieveFileBucketRequest) (*pb.RetrieveFileBucketResponse, error) {
	return c.service.RetrieveFileBucket(ctx, req)
}

// GetBucketAdmins returns the list of bucket admins (and optional owner) for the given bucket ID.
func (c *Client) GetBucketAdmins(ctx context.Context, req *pb.GetBucketAdminsRequest) (*pb.GetBucketAdminsResponse, error) {
	return c.service.GetBucketAdmins(ctx, req)
}

// IsBucketProtected returns whether the bucket requires a password (X-Bucket-Token) for access.
func (c *Client) IsBucketProtected(ctx context.Context, req *pb.IsBucketProtectedRequest) (*pb.IsBucketProtectedResponse, error) {
	return c.service.IsBucketProtected(ctx, req)
}

// AuthenticateBucket verifies the bucket password and returns a short-lived bucket access token.
func (c *Client) AuthenticateBucket(ctx context.Context, req *pb.AuthenticateBucketRequest) (*pb.AuthenticateBucketResponse, error) {
	return c.service.AuthenticateBucket(ctx, req)
}

// DeleteBucket deletes the bucket, its files in S3, and DB rows (files, bucket_admins). Returns files_deleted and error.
func (c *Client) DeleteBucket(ctx context.Context, req *pb.DeleteBucketRequest) (*pb.DeleteBucketResponse, error) {
	return c.service.DeleteBucket(ctx, req)
}
