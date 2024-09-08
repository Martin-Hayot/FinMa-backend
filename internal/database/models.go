package database

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID      `json:"id" gorm:"primary_key"`
	FirstName     string         `json:"first_name" validate:"required"`
	LastName      string         `json:"last_name" validate:"required"`
	Email         string         `json:"email" gorm:"uniqueIndex" validate:"required,email"`
	Password      string         `json:"password" validate:"required"`
	Role          string         `json:"role"`
	Transactions  []Transaction  `json:"transactions" gorm:"foreignKey:UserID"`
	BankAccounts  []BankAccount  `json:"bank_accounts" gorm:"foreignKey:UserID"`
	Budgets       []Budget       `json:"budgets" gorm:"foreignKey:UserID"`
	Notifications []Notification `json:"notifications" gorm:"foreignKey:UserID"`
	RefreshTokens []RefreshToken `json:"refresh_tokens" gorm:"foreignKey:UserID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type RefreshToken struct {
	ID        uuid.UUID `json:"id" gorm:"primary_key"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`

	UserID uuid.UUID `json:"user_id"`
	User   User      `json:"user"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BankAccount struct {
	ID            uuid.UUID `json:"id" gorm:"primary_key"`
	BankName      string    `json:"bank_name"`
	AccountType   string    `json:"account_type"`
	AccountNumber string    `json:"account_number" gorm:"uniqueIndex"`
	Balance       float64   `json:"balance"`

	UserID       uuid.UUID     `json:"user_id"`
	User         User          `json:"user"`
	Transactions []Transaction `json:"transactions" gorm:"foreignKey:BankAccountID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type Transaction struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key"`
	Category    string    `json:"category"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`
	IsRecuring  bool      `json:"is_recuring"`
	Description string    `json:"description"`

	UserID        uuid.UUID   `json:"user_id"`
	User          User        `json:"user"`
	BankAccountID uuid.UUID   `json:"bank_account_id"`
	BankAccount   BankAccount `json:"bank_account"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type Budget struct {
	ID        uuid.UUID `json:"id" gorm:"primary_key"`
	Category  string    `json:"category"`
	Amount    float64   `json:"amount"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`

	UserID uuid.UUID `json:"user_id"`
	User   User      `json:"user"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type Notification struct {
	ID       uuid.UUID `json:"id" gorm:"primary_key"`
	Type     string    `json:"type"`
	Message  string    `json:"message"`
	IsActive bool      `json:"is_active"`

	UserID uuid.UUID `json:"user_id"`
	User   User      `json:"user"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
