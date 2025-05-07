package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "POST, GET, OPTIONS, PUT, DELETE",
		AllowHeaders:     "Content-Type, Authorization, Accept, Origin",
		ExposeHeaders:    "Set-Cookie",
		AllowCredentials: true,
	}))

	// [Groups]
	// Frontend routes
	api := s.Group("/api")
	auth := api.Group("/auth")

	// [Routes]
	// General API routes
	api.Get("/", s.ShowAllRoutes)
	api.Get("/routes", s.ShowAllRoutes)
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

func (s *FiberServer) ShowAllRoutes(c *fiber.Ctx) error {
	routes := s.App.GetRoutes()
	// Filter out the routes that are not needed
	for i := len(routes) - 1; i >= 0; i-- {
		if routes[i].Path == "/api" || routes[i].Path == "/api/" {
			routes = append(routes[:i], routes[i+1:]...)
		}
	}
	return c.JSON(routes)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
