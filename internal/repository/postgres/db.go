package postgres

import (
	"fmt"

	"github.com/charmbracelet/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"FinMa/config"
	"FinMa/internal/domain"
)

// DB wraps the GORM database connection
type DB struct {
	DB *gorm.DB
}

// NewDB creates a new database connection
func NewDB(config *config.Config) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		config.Database.Host, config.Database.User, config.Database.Password, config.Database.Name, config.Database.Port, config.Database.SSLMode,
	)

	log.Info("Connecting to database", "host", config.Database.Host, "database", config.Database.Name)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(0)

	log.Info("Connected to database successfully")

	return &DB{DB: db}, nil
}

// Migrate runs database migrations
func (db *DB) Migrate() error {
	log.Info("Running database migrations")

	// List all models to auto-migrate
	err := db.DB.AutoMigrate(
		&domain.User{},
		&domain.Transaction{},
		&domain.Budget{},
		&domain.BankAccount{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Info("Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
