package server

import (
	"FinMa/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// [Groups]
	api := s.Group("/api")
	auth := api.Group("/auth")

	// [Middlewares]
	// api.Use(middlewares.Authorize)

	// [Routes]
	auth.Post("/signup", handlers.SignUpHandler)
	auth.Post("/login", handlers.LoginHandler)
	auth.Post("/refresh", handlers.RefreshHandler)
	api.Get("/", s.HelloWorldHandler)
	api.Get("/health", s.healthHandler)

}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
