package service

import (
	"fmt"
	"time"

	localpkg "github.com/cthulhu-platform/filemanager/internal/pkg"
	"github.com/cthulhu-platform/filemanager/pkg"
	"github.com/golang-jwt/jwt/v5"
)

// GenerateBucketAccessToken generates a JWT token for bucket access
func GenerateBucketAccessToken(bucketID string, userID *string, authTokenID *string, privileges []string) (string, error) {
	if bucketID == "" {
		return "", fmt.Errorf("bucket_id is required")
	}

	jwtSecret := localpkg.BUCKET_TOKEN_SECRET_KEY
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable is required")
	}

	now := time.Now()
	claims := &pkg.BucketAccessClaims{
		BucketID:    bucketID,
		Privileges:  privileges,
		UserID:      userID,
		AuthTokenID: authTokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(localpkg.BUCKET_TOKEN_EXPIRATION)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign bucket access token: %w", err)
	}

	return tokenString, nil
}

// ValidateBucketAccessToken validates a bucket access token and returns its claims
func ValidateBucketAccessToken(tokenString string) (*pkg.BucketAccessClaims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token is required")
	}

	jwtSecret := localpkg.BUCKET_TOKEN_SECRET_KEY
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	claims := &pkg.BucketAccessClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
