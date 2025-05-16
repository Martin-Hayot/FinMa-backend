package server

import (
	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// [Groups]
	api := s.Group("/api")
	auth := api.Group("/auth")

	// [Routes]
	// General API routes
	api.Get("/", s.ShowAllRoutes)
	api.Get("/routes", s.ShowAllRoutes)
	api.Get("/health", s.healthHandler)

	// [Plaid Routes]
	api.Post("/link/token/create", s.CreateLinkTokenHandler)

	// [Auth Routes]
	auth.Get("/me", s.Authorize("user"), s.MeHandler)
	auth.Post("/verify", s.Authorize("user"), s.VerifyHandler)
	auth.Post("/signup", s.SignUpHandler)
	auth.Post("/login", s.LoginHandler)
	auth.Post("/logout", s.LogoutHandler)
	auth.Post("/refresh", s.RefreshHandler)

	// Bank account routes
	api.Post("/bank-accounts", s.Authorize("user"), s.CreateBankAccount)

	// Transaction routes
	api.Post("/transactions", s.Authorize("user"), s.CreateTransaction)
	api.Get("/transactions", s.Authorize("user"), s.GetTransactions)
	api.Get("/transactions/:id", s.Authorize("user"), s.GetTransactionByID)
}

func (s *FiberServer) CreateLinkTokenHandler(c *fiber.Ctx) error {
	return nil
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
