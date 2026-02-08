package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	localPkg "github.com/cthulhu-platform/auth/internal/pkg"
	"github.com/cthulhu-platform/auth/pkg"
	"github.com/golang-jwt/jwt/v5"
)

// GenerateAccessToken creates a JWT with user claims.
func generateAccessToken(userID, email, provider string) (string, error) {
	secret := localPkg.JWT_SECRET_KEY
	claims := pkg.Claims{
		UserID:   userID,
		Email:    email,
		Provider: provider,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(localPkg.ACCESS_TOKEN_EXPIRATION)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(secret))
}

// validateAccessToken parses and verifies the JWT, returns claims or error.
func validateAccessToken(tokenString string) (*pkg.Claims, error) {
	secret := localPkg.JWT_SECRET_KEY
	tok, err := jwt.ParseWithClaims(tokenString, &pkg.Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(*pkg.Claims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

// generateRefreshToken returns a new random refresh token and its SHA256 hex hash for storage.
func generateRefreshToken() (plain, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}
	plain = hex.EncodeToString(b)
	h := sha256.Sum256([]byte(plain))
	hash = hex.EncodeToString(h[:])
	return plain, hash, nil
}
