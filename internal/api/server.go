package api

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"FinMa/config"
	"FinMa/internal/api/handlers"
	"FinMa/internal/repository/postgres"
	"FinMa/internal/service"
	"FinMa/plaid"
)

// Server represents the API server
type Server struct {
	app      *fiber.App
	config   *config.Config
	services *service.Services
	handlers *handlers.Handlers
}

// NewServer creates a new API server
func NewServer(config *config.Config, db *postgres.DB) *Server {
	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
		AppName:      "FinMa API",
	})

	// Setup global middleware
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://localhost:8080",
		AllowMethods:     "POST, GET, PATCH, OPTIONS, PUT, DELETE",
		AllowHeaders:     "Content-Type, Authorization, Accept, Origin, Access-Control-Allow-Origin",
		ExposeHeaders:    "Set-Cookie",
		AllowCredentials: true,
	}))

	// Create repositories
	userRepo := postgres.NewUserRepository(db.DB)
	bankAccountRepo := postgres.NewBankAccountRepository(db.DB)
	plaidItemRepo := postgres.NewPlaidItemRepository(db.DB)

	// Create validator service
	validatorService := service.NewValidatorService()

	// Create services
	authService := service.NewAuthService(userRepo, config)
	userService := service.NewUserService(userRepo)
	bankAccountService := service.NewBankAccountService(bankAccountRepo, plaidItemRepo, userRepo)
	plaidItemService := service.NewPlaidItemService(plaidItemRepo, userRepo)
	// Create services container
	services := &service.Services{
		Auth:        authService,
		User:        userService,
		BankAccount: bankAccountService,
		Validator:   validatorService,
		PlaidItem:   plaidItemService,
	}

	// Initialize Plaid client
	plaidClient := plaid.NewClient(config)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService, validatorService)
	userHandler := handlers.NewUserHandler(userService, validatorService)
	plaidHandler := handlers.NewPlaidHandler(*plaidClient, validatorService, userService, bankAccountService, plaidItemService)

	// Create handlers container
	handlers := &handlers.Handlers{
		Auth:  *authHandler,
		User:  *userHandler,
		Plaid: *plaidHandler,
	}

	// Create server
	server := &Server{
		app:      app,
		config:   config,
		services: services,
		handlers: handlers,
	}

	// Setup routes
	SetupRoutes(app, services, handlers)

	return server
}

// Start starts the server
func (s *Server) Start() {
	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%s", s.config.Port)
		log.Info("Starting server", "address", addr)

		if err := s.app.Listen(addr); err != nil {
			log.Fatal("Server error", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Gracefully shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.app.ShutdownWithContext(ctx); err != nil {
		log.Fatal("Server forced to shutdown", "error", err)
	}

	log.Info("Server stopped gracefully")
}

// Custom error handler
func customErrorHandler(c *fiber.Ctx, err error) error {
	// Default status code is 500
	code := fiber.StatusInternalServerError

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Log error
	log.Error("API error", "path", c.Path(), "error", err.Error())

	// Return JSON response
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
