package client

import (
	"context"
	"fmt"
	"time"

	"github.com/cthulhu-platform/lifecycle/pkg"
	"github.com/cthulhu-platform/lifecycle/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.LifecycleServiceClient
}

func NewClient(ctx context.Context, addr string) (*Client, error) {
	// TODO: add TLS credentials when we have a certificate
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	service := pb.NewLifecycleServiceClient(conn)

	return &Client{conn: conn, service: service}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) PostLifecycle(ctx context.Context, bucketSlug string, expiresAt time.Time) (*pkg.Lifecycle, error) {
	r, err := c.service.PostLifecycle(ctx, &pb.PostLifecycleRequest{
		BucketSlug: bucketSlug,
		ExpiresAt:  timestamppb.New(expiresAt),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to post lifecycle: %v", err)
	}

	return &pkg.Lifecycle{
		ID:         int(r.Lifecycle.Id),
		BucketSlug: r.Lifecycle.BucketSlug,
		ExpiresAt:  r.Lifecycle.ExpiresAt.AsTime(),
		CreatedAt:  r.Lifecycle.CreatedAt.AsTime(),
		UpdatedAt:  r.Lifecycle.UpdatedAt.AsTime(),
	}, nil
}

func (c *Client) GetLifecycle(ctx context.Context, bucketSlug string) (*pkg.Lifecycle, error) {
	r, err := c.service.GetLifecycle(ctx, &pb.GetLifecycleRequest{
		BucketSlug: bucketSlug,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get lifecycle: %v", err)
	}
	return &pkg.Lifecycle{
		ID:         int(r.Lifecycle.Id),
		BucketSlug: r.Lifecycle.BucketSlug,
		ExpiresAt:  r.Lifecycle.ExpiresAt.AsTime(),
		CreatedAt:  r.Lifecycle.CreatedAt.AsTime(),
		UpdatedAt:  r.Lifecycle.UpdatedAt.AsTime(),
	}, nil
}

func (c *Client) DeleteLifecycle(ctx context.Context, bucketSlug string) (bool, error) {
	r, err := c.service.DeleteLifecycle(ctx, &pb.DeleteLifecycleRequest{
		BucketSlug: bucketSlug,
	})
	if err != nil {
		return false, fmt.Errorf("failed to delete lifecycle: %v", err)
	}

	return r.Success, nil
}
