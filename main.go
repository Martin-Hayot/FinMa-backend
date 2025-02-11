package main

import (
	"FinMa/internal/server"
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	_ "github.com/joho/godotenv/autoload"
)

type Payload struct {
	UserID   int
	Username string
	Email    string
}

func main() {
	log.SetLevel(log.DebugLevel)
	server := server.New()

	server.RegisterFiberRoutes()
	server.Use(helmet.New())
	server.Use(limiter.New())
	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := server.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
