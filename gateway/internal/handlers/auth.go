package handlers

import (
	"fmt"

	"github.com/cthulhu-platform/gateway/internal/connections"
	"github.com/cthulhu-platform/gateway/internal/pkg"
	"github.com/gofiber/fiber/v2"
)

func OAuthInitiate(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		provider := c.Params("provider")
		if provider == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "provider parameter is required",
			})
		}

		oauthURL, err := conns.Auth.InitiateOAuth(c.Context(), provider)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Redirect(oauthURL, fiber.StatusFound)
	}
}

func OAuthCallback(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		provider := c.Params("provider")
		code := c.Query("code")
		state := c.Query("state")

		if provider == "" || code == "" || state == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "provider, code, and state are required",
			})
		}

		acceptHeader := c.Get("Accept")
		isAPIRequest := acceptHeader == "application/json" || c.Get("Content-Type") == "application/json"

		if isAPIRequest {
			authResponse, err := conns.Auth.HandleOAuthCallback(c.Context(), provider, code, state)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			return c.JSON(authResponse)
		}

		clientCallbackURL := pkg.CORS_ORIGIN + "/signin/callback"
		redirectURL := fmt.Sprintf("%s?code=%s&state=%s", clientCallbackURL, code, state)
		return c.Redirect(redirectURL, fiber.StatusFound)
	}
}

func TokenRefresh(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		if req.RefreshToken == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "refresh_token is required",
			})
		}

		tokenPair, err := conns.Auth.RefreshToken(c.Context(), req.RefreshToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(tokenPair)
	}
}

func TokenLogout(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "authorization header is required",
			})
		}

		var accessToken string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			accessToken = authHeader[7:]
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}

		_, err := conns.Auth.Logout(c.Context(), accessToken)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "logged out successfully",
		})
	}
}

func TokenValidate(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "authorization header is required",
			})
		}

		var token string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}

		user, err := conns.Auth.ValidateToken(c.Context(), token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"valid": true,
			"user":  user,
		})
	}
}
