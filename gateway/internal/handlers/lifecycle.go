package handlers

import (
	"strings"
	"time"

	"github.com/cthulhu-platform/gateway/internal/connections"
	"github.com/gofiber/fiber/v2"
)

func GetNormalizedBucketLifecycle(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := strings.TrimSpace(c.Params("id"))
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bucket id is required"})
		}

		lifecycle, err := conns.Lifecycle.GetLifecycle(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "lifecycle not found"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"bucket_id":  lifecycle.BucketSlug,
			"expires_at": lifecycle.ExpiresAt.UTC().Format(time.RFC3339),
		})
	}
}
