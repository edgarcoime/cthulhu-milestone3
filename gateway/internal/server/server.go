package server

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/cthulhu-platform/gateway/internal/connections"
	"github.com/cthulhu-platform/gateway/internal/pkg"
	"github.com/cthulhu-platform/gateway/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"
)

type FiberServerConfig struct {
	Host   string
	Port   string
	Logger *slog.Logger
}

type FiberServer struct {
	Host   string
	Port   string
	Conns  *connections.ConnectionsContainer
	Logger *slog.Logger
}

func NewFiberServer(cfg FiberServerConfig, conns *connections.ConnectionsContainer) *FiberServer {
	return &FiberServer{
		Host:   cfg.Host,
		Port:   cfg.Port,
		Logger: cfg.Logger,
		Conns:  conns,
	}
}

func (s *FiberServer) Start() {
	// Setup dependencies
	app := fiber.New(fiber.Config{
		BodyLimit: pkg.BODY_LIMIT_MB * 1024 * 1024,
	})

	// Setup middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: pkg.CORS_ORIGIN,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Bucket-Token",
	}))
	slogCfg := slogfiber.Config{
		WithClientIP: true,
	}
	app.Use(slogfiber.NewWithConfig(s.Logger, slogCfg))
	app.Use(recover.New())

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from the server!")
	})
	routes.TestingRouter(app, s.Conns)
	routes.FilesRouter(app, s.Conns)
	routes.LifecycleRouter(app, s.Conns)
	routes.AuthRouter(app, s.Conns)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		slog.Info("Shutting down gracefully...")
		app.Shutdown()
	}()

	if err := app.Listen(s.Host + ":" + s.Port); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
