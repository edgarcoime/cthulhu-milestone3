package routes

import (
	"github.com/cthulhu-platform/gateway/internal/connections"
	"github.com/cthulhu-platform/gateway/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func AuthRouter(app fiber.Router, conns *connections.ConnectionsContainer) {
	// OAuth endpoints
	app.Get("/auth/oauth/:provider", handlers.OAuthInitiate(conns))
	app.Get("/auth/oauth/:provider/callback", handlers.OAuthCallback(conns))

	// Token management
	app.Post("/auth/refresh", handlers.TokenRefresh(conns))
	app.Post("/auth/logout", handlers.TokenLogout(conns))
	app.Post("/auth/validate", handlers.TokenValidate(conns))
}
