package main

import (
	"FinMa/config"
	"FinMa/internal/server"

	"github.com/charmbracelet/log"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	log.SetLevel(log.DebugLevel)
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize configuration
	cfg := config.LoadConfig()
	server := server.New(cfg)

	log.Fatal(server.Listen(":" + cfg.Port))
}
