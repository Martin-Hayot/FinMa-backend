package server

import (
	"FinMa/config"
	"FinMa/internal/database"
	"FinMa/plaid"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type FiberServer struct {
	*fiber.App
	cfg         *config.Config
	db          database.Service
	plaidClient *plaid.Client
}

func New(cfg *config.Config) *FiberServer {

	// Create Plaid client
	plaidClient := plaid.NewClient(cfg)

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "FinMa",
			AppName:      "FinMa",
		}),
		db:          database.New(),
		cfg:         cfg,
		plaidClient: plaidClient,
	}
	server.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://localhost:8080",
		AllowMethods:     "POST, GET, OPTIONS, PUT, DELETE",
		AllowHeaders:     "Content-Type, Authorization, Accept, Origin, Access-Control-Allow-Origin",
		ExposeHeaders:    "Set-Cookie",
		AllowCredentials: true,
	}))

	server.RegisterFiberRoutes()

	return server
}
