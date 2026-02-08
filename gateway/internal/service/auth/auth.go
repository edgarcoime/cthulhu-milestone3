package auth

import "context"

type AuthService interface {
	InitiateOAuth(provider string) (string, error)
	HandleOAuthCallback(provider string, code string) (string, error)
}

type authServiceRMQ struct {
	// Connections to rabbitmq here
}

func NewAuthServiceRMQ(ctx context.Context) AuthService {
	_ = ctx // use for connection init with timeout
	return &authServiceRMQ{
		// Connections to rabbitmq here
	}
}

func (s *authServiceRMQ) InitiateOAuth(provider string) (string, error) {
	// TODO: complete RPC logic
	return "", nil
}

func (s *authServiceRMQ) HandleOAuthCallback(provider string, code string) (string, error) {
	// TODO: complete RPC logic
	return "", nil
}
