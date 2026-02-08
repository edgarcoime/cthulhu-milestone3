package routes

import (
	"github.com/cthulhu-platform/gateway/internal/connections"
	"github.com/cthulhu-platform/gateway/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func LifecycleRouter(app fiber.Router, conns *connections.ConnectionsContainer) {
	app.Get("/lifecycle/s/:id", handlers.GetNormalizedBucketLifecycle(conns))
}
