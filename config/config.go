package config

import (
	"log"
	"os"
	"strings"

	plaid "github.com/plaid/plaid-go/v31/plaid"
)

// Config holds all configuration for the application
type Config struct {
	PlaidClientID     string
	PlaidSecret       string
	PlaidEnv          plaid.Environment
	PlaidProducts     []string
	PlaidCountryCodes []string
	PlaidRedirectURI  string
	Port              string
	DatabaseHost      string
	DatabasePort      string
	DatabaseUser      string
	DatabasePassword  string
	DatabaseName      string
	DatabaseSchema    string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		PlaidClientID:     getEnv("PLAID_CLIENT_ID", ""),
		PlaidSecret:       getEnv("PLAID_SECRET", ""),
		PlaidProducts:     strings.Split(getEnv("PLAID_PRODUCTS", "transactions"), ","),
		PlaidCountryCodes: strings.Split(getEnv("PLAID_COUNTRY_CODES", "US"), ","),
		PlaidRedirectURI:  getEnv("PLAID_REDIRECT_URI", ""),
		Port:              getEnv("PORT", "8080"),
		DatabaseHost:      getEnv("DB_HOST", "localhost"),
		DatabasePort:      getEnv("DB_PORT", "5432"),
		DatabaseUser:      getEnv("DB_USERNAME", "postgres"),
		DatabasePassword:  getEnv("DB_PASSWORD", "postgres"),
		DatabaseName:      getEnv("DB_DATABASE", "finma"),
		DatabaseSchema:    getEnv("DB_SCHEMA", "public"),
	}

	// Validate required configuration
	if config.PlaidClientID == "" || config.PlaidSecret == "" || config.DatabaseUser == "" || config.DatabasePassword == "" {
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
