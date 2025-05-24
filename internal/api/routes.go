package api

import (
	"FinMa/internal/api/handlers"
	"FinMa/internal/api/middleware"
	"FinMa/internal/service"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(app *fiber.App, services *service.Services, handlers *handlers.Handlers) {
	// API group
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/signup", handlers.Auth.SignUp)
	auth.Post("/login", handlers.Auth.Login)
	auth.Post("/refresh", handlers.Auth.Refresh)
	auth.Post("/logout", handlers.Auth.Logout)

	// Protected routes
	protected := api.Group("", middleware.AuthMiddleware(services.Auth))
	protected.Get("/me", handlers.Auth.Me)

	// User routes
	users := protected.Group("/users")
	users.Patch("/:id", handlers.User.Update)

	// Plaid routes
	plaid := protected.Group("/plaid")
	plaid.Post("/create_link_token", handlers.Plaid.CreateLinkToken)
	plaid.Post("/exchange_public_token", handlers.Plaid.ExchangePublicToken)

	// Transaction routes
	// transactions := protected.Group("/transactions")
	// transactions.Get("/", handlers.Transaction.GetAll)
	// transactions.Post("/", handlers.Transaction.Create)
	// transactions.Get("/:id", handlers.Transaction.GetByID)
	// transactions.Put("/:id", handlers.Transaction.Update)
	// transactions.Delete("/:id", handlers.Transaction.Delete)

	// Budget routes
	// budgets := protected.Group("/budgets")
	// budgets.Get("/", handlers.Budget.GetAll)
	// budgets.Post("/", handlers.Budget.Create)
	// budgets.Get("/:id", handlers.Budget.GetByID)
	// budgets.Put("/:id", handlers.Budget.Update)
	// budgets.Delete("/:id", handlers.Budget.Delete)

	// Settings routes
	// settings := protected.Group("/settings")
	// settings.Get("/", handlers.Settings.Get)
	// settings.Put("/", handlers.Settings.Update)
}
