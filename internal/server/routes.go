package server

import (
	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// [Groups]
	api := s.Group("/api")
	auth := api.Group("/auth")

	// [Middlewares]
	// api.Use(middlewares.Authorize)

	// [Routes]
	// General routes
	api.Get("/", s.HelloWorldHandler)
	api.Get("/health", s.healthHandler)

	// Auth routes
	auth.Post("/signup", s.SignUpHandler)
	auth.Post("/login", s.LoginHandler)
	auth.Post("/refresh", s.RefreshHandler)

	// Bank account routes
	api.Post("/bank-accounts", s.Authorize("user"), s.CreateBankAccount)

	// Transaction routes
	api.Post("/transactions", s.Authorize("user"), s.CreateTransaction)
	api.Get("/transactions", s.Authorize("user"), s.GetTransactions)
	api.Get("/transactions/:id", s.Authorize("user"), s.GetTransactionByID)

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
