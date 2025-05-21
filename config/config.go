package config

import (
	"log"
	"os"
	"strings"

	plaid "github.com/plaid/plaid-go/v31/plaid"
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

type Config struct {
	PlaidClientID      string
	PlaidSecret        string
	PlaidEnv           plaid.Environment
	PlaidProducts      []string
	PlaidCountryCodes  []string
	PlaidRedirectURI   string
	Port               string
	AccessTokenSecret  string
	RefreshTokenSecret string
	Database           DatabaseConfig
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		AccessTokenSecret:  getEnv("ACCESS_TOKEN_SECRET", "default_secret"),
		RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", "default_secret"),
		Port:               getEnv("PORT", "8080"),
		PlaidClientID:      getEnv("PLAID_CLIENT_ID", ""),
		PlaidSecret:        getEnv("PLAID_SECRET", ""),
		PlaidProducts:      strings.Split(getEnv("PLAID_PRODUCTS", "transactions"), ","),
		PlaidCountryCodes:  strings.Split(getEnv("PLAID_COUNTRY_CODES", "US"), ","),
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
	if config.PlaidClientID == "" || config.PlaidSecret == "" || config.Database.User == "" || config.Database.Password == "" {
		log.Fatal("Error: PLAID_SECRET, PLAID_CLIENT_ID, DB_USERNAME or DB_PASSWORD is not set. Did you copy .env.example to .env and fill it out?")
	}

	// Set Plaid environment
	plaidEnv := getEnv("PLAID_ENV", "sandbox")
	if plaidEnv == "production" {
		config.PlaidEnv = plaid.Production
	} else {
		config.PlaidEnv = plaid.Sandbox
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
