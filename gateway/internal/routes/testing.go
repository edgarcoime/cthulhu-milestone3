package routes

import (
	"github.com/cthulhu-platform/gateway/internal/connections"
	"github.com/cthulhu-platform/gateway/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func TestingRouter(app fiber.Router, conns *connections.ConnectionsContainer) {
	app.Get("/testing/fanout", handlers.FanoutPingAllServices(conns))
	app.Get("/diagnose/filemanager", handlers.DiagnoseFilemanagerService(conns))
	app.Get("/diagnose/auth", handlers.DiagnoseAuthService(conns))
	app.Get("/diagnose/lifecycle", handlers.DiagnoseLifecycleService(conns))
	app.Get("/diagnose/lifecycle/delete", handlers.DiagnoseLifecycleServiceDelete(conns))
}
