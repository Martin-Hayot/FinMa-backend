package server

import (
	"github.com/gofiber/fiber/v2"

	"FinMa/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "FinMa",
			AppName:      "FinMa",
		}),

		db: database.New(),
	}

	return server
}
