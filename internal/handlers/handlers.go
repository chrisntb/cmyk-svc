package handlers

import (
	"cmyk/internal/clients/env"
	"cmyk/internal/clients/k8s"
	"cmyk/internal/clients/mock"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Actions struct for database DML
type Handlers struct {
	App        *fiber.App
	EnvClient  *env.Client
	K8sClient  *k8s.Client
	MockClient *mock.Client
}

func NewHandlers(app *fiber.App, envClient *env.Client, k8sClient *k8s.Client, mockClient *mock.Client) Handlers {
	handlers := Handlers{App: app, EnvClient: envClient, K8sClient: k8sClient, MockClient: mockClient}

	// Middleware
	handlers.App.Use(recover.New())
	handlers.App.Use(logger.New())
	handlers.App.Use(cors.New())

	// Static files
	handlers.App.Static("/static", "./static")

	// API routes
	root := handlers.App.Group("/api")
	root.Get("/health", handlers.Health)

	v1 := handlers.App.Group("/api/v1")
	v1.Get("/nodes", handlers.ReadNodes)
	v1.Get("/nodes/:name", handlers.ReadNodeDetail)

	v1.Get("/pods", handlers.ReadPods)
	v1.Get("/namespaces/:namespace/pods/:name", handlers.ReadPodDetail)

	v1.Post("/jobs", handlers.CreateJob)

	v1.Get("/resource-flavors", handlers.ReadResourceFlavors)
	v1.Get("/resource-flavors/:name", handlers.ReadResourceFlavorDetail)
	v1.Post("/resource-flavors", handlers.CreateResourceFlavor)
	v1.Delete("/resource-flavors/:name", handlers.DeleteResourceFlavor)

	v1.Get("/local-queues", handlers.ReadLocalQueues)
	v1.Get("/namespaces/:namespace/local-queues/:name", handlers.ReadLocalQueueDetail)
	v1.Post("/namespaces/:namespace/local-queues", handlers.CreateLocalQueue)
	v1.Delete("/namespaces/:namespace/local-queues/:name", handlers.DeleteLocalQueue)

	v1.Get("/kai-scheduler-queues", handlers.ReadKaiSchedulerQueues)
	v1.Get("/kai-scheduler-queues/:name/child-queues", handlers.ReadKaiSchedulerChildQueues)

	// Must come last
	handlers.App.Use(NotFound)

	return handlers
}

// NotFound returns custom 404 page
func NotFound(c *fiber.Ctx) error {
	return c.Status(404).SendFile("./static/private/404.html")
}
