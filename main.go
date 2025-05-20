package main

import (
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"

	"FinMa/config"
	"FinMa/internal/api"
	"FinMa/internal/repository/postgres"
)

func main() {
	log.SetLevel(log.DebugLevel)
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize configuration
	cfg := config.LoadConfig()

	// Connect to database
	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}

	// Run database migrations
	if err := db.Migrate(); err != nil {
		log.Fatal("Failed to run database migrations", "error", err)
	}

	server := api.NewServer(cfg, db)

	server.Start()
}
