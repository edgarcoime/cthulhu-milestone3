package middleware

import (
	"strings"

	"github.com/cthulhu-platform/auth/pkg"
	"github.com/cthulhu-platform/gateway/internal/connections"
	fmpb "github.com/cthulhu-platform/proto/pkg/filemanager"
	"github.com/gofiber/fiber/v2"
)

const (
	LocalsKeyUserID = "user_id"
	LocalsKeyUser   = "user"
)

// RequireAuth validates the Bearer token via the auth service and attaches user to context.
// Returns 401 if the token is missing, malformed, or invalid.
func RequireAuth(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "authorization header is required"})
		}
		if len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid authorization header format"})
		}
		token := authHeader[7:]

		user, err := conns.Auth.ValidateToken(c.Context(), token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}

		c.Locals(LocalsKeyUserID, user.ID)
		c.Locals(LocalsKeyUser, user)
		return c.Next()
	}
}

// OptionalAuth validates the Bearer token if present and attaches user to context.
// If the header is missing or the token is invalid, the request continues without user in context.
func OptionalAuth(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}
		if len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
			return c.Next()
		}
		token := authHeader[7:]

		user, err := conns.Auth.ValidateToken(c.Context(), token)
		if err != nil {
			return c.Next()
		}

		c.Locals(LocalsKeyUserID, user.ID)
		c.Locals(LocalsKeyUser, user)
		return c.Next()
	}
}

// GetUser returns the authenticated user from context, or nil if not set.
func GetUser(c *fiber.Ctx) *pkg.UserInfo {
	v := c.Locals(LocalsKeyUser)
	if v == nil {
		return nil
	}
	u, _ := v.(*pkg.UserInfo)
	return u
}

// BucketAuth runs optional JWT validation (sets user if Bearer valid), then for the bucket in :id
// calls filemanager IsBucketProtected; if protected and X-Bucket-Token is missing returns 401.
func BucketAuth(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token := authHeader[7:]
			if user, err := conns.Auth.ValidateToken(c.Context(), token); err == nil {
				c.Locals(LocalsKeyUserID, user.ID)
				c.Locals(LocalsKeyUser, user)
			}
		}
		bucketID := strings.TrimSpace(c.Params("id"))
		if bucketID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bucket id is required"})
		}
		res, err := conns.Filemanager.IsBucketProtected(c.Context(), &fmpb.IsBucketProtectedRequest{BucketId: bucketID})
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": err.Error()})
		}
		if res != nil && res.Error != "" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": res.Error})
		}
		if res != nil && res.Protected && c.Get("X-Bucket-Token") == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bucket is protected; X-Bucket-Token is required"})
		}
		return c.Next()
	}
}
