package client

import (
	"context"
	"fmt"

	"github.com/cthulhu-platform/auth/pkg"
	pb "github.com/cthulhu-platform/proto/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AuthServiceClient
}

func NewClient(ctx context.Context, addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	service := pb.NewAuthServiceClient(conn)

	return &Client{conn: conn, service: service}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) InitiateOAuth(ctx context.Context, provider string) (string, error) {
	r, err := c.service.InitiateOAuth(ctx, &pb.InitiateOAuthRequest{Provider: provider})
	if err != nil {
		return "", fmt.Errorf("failed to initiate OAuth: %v", err)
	}
	return r.RedirectUrl, nil
}

func (c *Client) HandleOAuthCallback(ctx context.Context, provider string, code string, state string) (*pkg.AuthResponse, error) {
	r, err := c.service.HandleOAuthCallback(ctx, &pb.HandleOAuthCallbackRequest{Provider: provider, Code: code, State: state})
	if err != nil {
		return nil, fmt.Errorf("failed to handle OAuth callback: %v", err)
	}
	return &pkg.AuthResponse{AccessToken: r.AccessToken, RefreshToken: r.RefreshToken, User: &pkg.UserInfo{ID: r.User.Id, Email: r.User.Email, Username: r.User.Username, AvatarUrl: r.User.AvatarUrl}}, nil
}

func (c *Client) ValidateToken(ctx context.Context, token string) (*pkg.UserInfo, error) {
	r, err := c.service.ValidateToken(ctx, &pb.ValidateTokenRequest{Token: token})
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}
	return &pkg.UserInfo{ID: r.User.Id, Email: r.User.Email, Username: r.User.Username, AvatarUrl: r.User.AvatarUrl}, nil
}

func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*pkg.TokenPair, error) {
	r, err := c.service.RefreshToken(ctx, &pb.RefreshTokenRequest{RefreshToken: refreshToken})
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %v", err)
	}
	return &pkg.TokenPair{AccessToken: r.AccessToken, RefreshToken: r.RefreshToken}, nil
}

func (c *Client) Logout(ctx context.Context, accessToken string) (bool, error) {
	r, err := c.service.Logout(ctx, &pb.LogoutRequest{AccessToken: accessToken})
	if err != nil {
		return false, fmt.Errorf("failed to logout: %v", err)
	}
	return r.Success, nil
}
