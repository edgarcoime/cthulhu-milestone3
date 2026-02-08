package handlers

import (
	"log"
	"time"

	"github.com/cthulhu-platform/gateway/internal/connections"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func FanoutPingAllServices(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"transaction_id": "TODO",
			"message":        "Pinging all services in network",
		})
	}
}

func DiagnoseFilemanagerService(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"transaction_id": "TODO",
			"message":        "Diagnosing filemanager service",
		})
	}
}

func DiagnoseAuthService(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"transaction_id": "TODO",
			"message":        "Diagnosing auth service",
		})
	}
}

func DiagnoseLifecycleServiceDelete(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		slug := c.Query("slug")
		if slug == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"transaction_id": "TODO",
				"message":        "Missing query parameter: slug",
			})
		}
		ok, err := conns.Lifecycle.DeleteLifecycle(ctx, slug)
		if err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"transaction_id": "TODO",
				"bucket_slug":    slug,
				"message":        "Failed to delete lifecycle",
				"error":          err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"transaction_id": "TODO",
			"bucket_slug":    slug,
			"success":        ok,
			"message":        "Diagnosing lifecycle service delete: deleted",
		})
	}
}

func DiagnoseLifecycleService(conns *connections.ConnectionsContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		uuid := uuid.New().String()
		created, err := conns.Lifecycle.PostLifecycle(ctx, uuid, time.Now().Add(1*time.Hour))
		if err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"transaction_id": "TODO",
				"bucket_slug":    uuid,
				"message":        "Failed to post lifecycle",
				"error":          err.Error(),
			})
		}

		log.Printf("Created lifecycle: %+v\n", created)
		// Sleep for 1 second to allow for write to committed
		// time.Sleep(1 * time.Second)

		lifecycle, err := conns.Lifecycle.GetLifecycle(ctx, uuid)
		if err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"transaction_id": "TODO",
				"bucket_slug":    uuid,
				"message":        "Failed to get lifecycle",
				"error":          err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"transaction_id": "TODO",
			"bucket_slug":    uuid,
			"lifecycle":      lifecycle,
			"message":        "Diagnosing lifecycle service: Successfully created and retrieved lifecycle",
		})
	}
}
