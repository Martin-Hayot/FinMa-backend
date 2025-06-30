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

	// GoCardless routes
	gocardless := protected.Group("/gocardless")
	gocardless.Get("/institutions/:country_code", handlers.GoCardless.GetInstitutions)
	gocardless.Post("/link", handlers.GoCardless.LinkAccount)
	gocardless.Patch("/requisitions/:id", handlers.GoCardless.SyncRequisition)
	gocardless.Get("/token/status", handlers.GoCardless.GetTokenStatus)

	// Bank Account routes
	bankAccounts := protected.Group("/bank-accounts")
	bankAccounts.Get("/", handlers.BankAccount.GetAccounts)

	accounts := protected.Group("/accounts")
	accounts.Get("/", handlers.BankAccount.GetAccounts)
	// accounts.Get("/:id", handlers.BankAccount.GetAccountDetails)
	// accounts.Get("/:id/balances", handlers.BankAccount.GetAccountBalances)
	// accounts.Get("/:id/transactions", handlers.BankAccount.GetAccountTransactions)
}
