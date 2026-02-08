package routes

import (
	"github.com/cthulhu-platform/gateway/internal/connections"
	"github.com/cthulhu-platform/gateway/internal/handlers"
	"github.com/cthulhu-platform/gateway/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func FilesRouter(app fiber.Router, conns *connections.ConnectionsContainer) {
	app.Post("/files/upload", middleware.OptionalAuth(conns), handlers.FileUpload(conns))
	app.Post("/files/upload/prepare", middleware.OptionalAuth(conns), handlers.FileUploadPrepare(conns))
	app.Post("/files/upload/confirm", middleware.OptionalAuth(conns), handlers.FileUploadConfirm(conns))
	app.Post("/files/s/:id/authenticate", middleware.OptionalAuth(conns), handlers.FileAuthenticate(conns))
	app.Get("/files/s/:id", middleware.BucketAuth(conns), handlers.FileBucketGet(conns))
	app.Get("/files/s/:id/admins", middleware.BucketAuth(conns), handlers.FileAdmins(conns))
	app.Get("/files/s/:id/protected", handlers.FileBucketProtected(conns))
	app.Get("/files/s/:id/d/:filename", middleware.BucketAuth(conns), handlers.FileDownload(conns))
}
