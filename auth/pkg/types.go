package pkg

import "github.com/golang-jwt/jwt/v5"

// TODO: where should this be? Generated from proto? or in pkg?
type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	User         *UserInfo `json:"user"`
}

type UserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Provider string `json:"provider"`
	jwt.RegisteredClaims
}
