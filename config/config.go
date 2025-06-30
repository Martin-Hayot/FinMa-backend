package config

import (
	"log"
	"os"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Schema   string
	SSLMode  string
}

type GoCardlessConfig struct {
	RedirectURL string
	ClientID    string
	Secret      string
}

type Config struct {
	Port               string
	AccessTokenSecret  string
	RefreshTokenSecret string
	Database           DatabaseConfig
	GoCardless         GoCardlessConfig
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		AccessTokenSecret:  getEnv("ACCESS_TOKEN_SECRET", "default_secret"),
		RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", "default_secret"),
		Port:               getEnv("PORT", "8080"),
		GoCardless: GoCardlessConfig{
			RedirectURL: getEnv("GOCARDLESS_REDIRECT_URL", "http://localhost:3000/gocardless/callback"),
			ClientID:    getEnv("GOCARDLESS_CLIENT_ID", ""),
			Secret:      getEnv("GOCARDLESS_SECRET", ""),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USERNAME", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_DATABASE", "finma"),
			Schema:   getEnv("DB_SCHEMA", "public"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
	// Set default values
	if config.GoCardless.ClientID == "" || config.GoCardless.Secret == "" || config.Database.User == "" || config.Database.Password == "" {
		log.Fatal("Error: GOCARDLESS_CLIENT_ID, GOCARDLESS_SECRET, DB_USERNAME or DB_PASSWORD is not set. Did you copy .env.example to .env and fill it out?")
	}

	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
