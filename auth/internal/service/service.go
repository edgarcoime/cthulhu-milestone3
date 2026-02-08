package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	localPkg "github.com/cthulhu-platform/auth/internal/pkg"
	"github.com/cthulhu-platform/auth/internal/repository"
	"github.com/cthulhu-platform/auth/internal/repository/sqlc/db"
	"github.com/cthulhu-platform/auth/pkg"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type Service interface {
	InitiateOAuth(ctx context.Context, provider string) (string, error)
	HandleOAuthCallback(ctx context.Context, provider string, code string, state string) (*pkg.AuthResponse, error)
	ValidateToken(ctx context.Context, token string) (*pkg.UserInfo, error)
	RefreshToken(ctx context.Context, refreshToken string) (*pkg.TokenPair, error)
	Logout(ctx context.Context, accessToken string) error
}

type authService struct {
	repo repository.Repository
}

func NewAuthService(repo repository.Repository) Service {
	return &authService{repo: repo}
}

func (s *authService) InitiateOAuth(ctx context.Context, provider string) (string, error) {
	// Generate PKCE values
	codeVerifier := generateCodeVerifier()
	codeChallenge := generateCodeChallenge(codeVerifier)

	// generate state
	state := uuid.New().String()

	now := time.Now()
	expiresAt := now.Add(localPkg.OAUTH_SESSION_EXPIRATION_TIME).Unix()

	// Store OAuth session (validated in HandleOAuthCallback via state)
	session := &db.OauthSession{
		State:         state,
		Provider:      provider,
		CodeVerifier:  codeVerifier,
		CodeChallenge: codeChallenge,
		RedirectUri:   getRedirectURI(provider),
		ExpiresAt:     expiresAt,
		CreatedAt:     now.Unix(),
	}

	if err := s.repo.CreateOAuthSession(ctx, session); err != nil {
		return "", err
	}

	return buildOAuthURL(provider, state, codeChallenge), nil
}

func (s *authService) HandleOAuthCallback(ctx context.Context, provider string, code string, state string) (*pkg.AuthResponse, error) {
	session, err := s.repo.GetOAuthSession(ctx, state)
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth session: %w", err)
	}
	if session == nil {
		return nil, fmt.Errorf("invalid or expired OAuth session")
	}
	if time.Now().Unix() > session.ExpiresAt {
		_ = s.repo.DeleteOAuthSession(ctx, state)
		return nil, fmt.Errorf("OAuth session expired")
	}
	if session.Provider != provider {
		return nil, fmt.Errorf("provider mismatch")
	}

	oauthConfig := getOAuthConfig(provider)
	if oauthConfig == nil {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
	token, err := oauthConfig.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", session.CodeVerifier))
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	userInfo, err := fetchUserInfo(provider, token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	existingUser, err := s.repo.GetUserByOAuthID(ctx, provider, userInfo.OAuthUserID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	now := time.Now().Unix()
	var user *db.User

	if existingUser != nil {
		if err := s.repo.UpdateUser(ctx, &db.User{
			ID:        existingUser.ID,
			Username:  ptrToNullString(userInfo.Username),
			AvatarUrl: ptrToNullString(userInfo.AvatarURL),
			UpdatedAt: now,
		}); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
		user = existingUser
		user.Username = ptrToNullString(userInfo.Username)
		user.AvatarUrl = ptrToNullString(userInfo.AvatarURL)
		user.UpdatedAt = now
	} else {
		user = &db.User{
			ID:            uuid.New().String(),
			OauthProvider: provider,
			OauthUserID:   userInfo.OAuthUserID,
			Email:         userInfo.Email,
			Username:      ptrToNullString(userInfo.Username),
			AvatarUrl:     ptrToNullString(userInfo.AvatarURL),
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := s.repo.CreateUser(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	accessToken, err := generateAccessToken(user.ID, user.Email, user.OauthProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshPlain, refreshHash, err := generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshExpiresAt := time.Now().Add(localPkg.REFRESH_TOKEN_EXPIRATION).Unix()
	if err := s.repo.CreateRefreshToken(ctx, &db.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: refreshExpiresAt,
		CreatedAt: now,
	}); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	_ = s.repo.DeleteOAuthSession(ctx, state)

	return &pkg.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshPlain,
		User:         userToUserInfo(user),
	}, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*pkg.UserInfo, error) {
	claims, err := validateAccessToken(token)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return userToUserInfo(user), nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*pkg.TokenPair, error) {
	hash := sha256Hex(refreshToken)

	tokenRecord, err := s.repo.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	if tokenRecord.RevokedAt.Valid {
		return nil, fmt.Errorf("refresh token has been revoked")
	}
	if time.Now().Unix() > tokenRecord.ExpiresAt {
		_ = s.repo.RevokeRefreshToken(ctx, tokenRecord.ID, "expired")
		return nil, fmt.Errorf("refresh token expired")
	}

	user, err := s.repo.GetUserByID(ctx, tokenRecord.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := s.repo.RevokeRefreshToken(ctx, tokenRecord.ID, "token_refreshed"); err != nil {
		return nil, fmt.Errorf("failed to revoke old token: %w", err)
	}

	accessToken, err := generateAccessToken(user.ID, user.Email, user.OauthProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshPlain, refreshHash, err := generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	now := time.Now().Unix()
	refreshExpiresAt := time.Now().Add(localPkg.REFRESH_TOKEN_EXPIRATION).Unix()
	if err := s.repo.CreateRefreshToken(ctx, &db.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: refreshExpiresAt,
		CreatedAt: now,
	}); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &pkg.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshPlain,
	}, nil
}

func (s *authService) Logout(ctx context.Context, accessToken string) error {
	claims, err := validateAccessToken(accessToken)
	if err != nil {
		return err
	}
	return s.repo.RevokeAllUserTokens(ctx, claims.UserID, "user_logout")
}

func ptrToNullString(s *string) sql.NullString {
	if s == nil || *s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

func userToUserInfo(u *db.User) *pkg.UserInfo {
	if u == nil {
		return nil
	}
	info := &pkg.UserInfo{ID: u.ID, Email: u.Email}
	if u.Username.Valid {
		info.Username = u.Username.String
	}
	if u.AvatarUrl.Valid {
		info.AvatarUrl = u.AvatarUrl.String
	}
	return info
}

func sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
