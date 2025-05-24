package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"primary_key"`
	FirstName  string    `validate:"required"`
	LastName   string    `validate:"required"`
	Email      string    `gorm:"uniqueIndex" validate:"required,email"`
	Password   string    `validate:"required"`
	Role       string
	IsVerified bool `gorm:"default:false"`

	// Associations
	PlaidItems    []PlaidItem    `gorm:"foreignKey:UserID"`
	BankAccounts  []BankAccount  `gorm:"foreignKey:UserID"`
	Transactions  []Transaction  `gorm:"foreignKey:UserID"`
	Budgets       []Budget       `gorm:"foreignKey:UserID"`
	Notifications []Notification `gorm:"foreignKey:UserID"`
	RefreshTokens []RefreshToken `gorm:"foreignKey:UserID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type PlaidItem struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	ItemID      string    `gorm:"uniqueIndex"` // From Plaid
	AccessToken string    `gorm:"uniqueIndex"` // From Plaid
	Institution string    // Optional: name of the bank (from Plaid metadata)

	UserID uuid.UUID
	User   User

	CreatedAt time.Time
	UpdatedAt time.Time
}

type EmailVerificationToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Token     string    `gorm:"not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
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
	ID               uuid.UUID `gorm:"primary_key"`
	AccountID        string    `gorm:"uniqueIndex"` // From Plaid
	AccountName      string
	OfficialName     string
	AccountType      string
	Mask             string // Last 4 digits etc
	CurrentBalance   float64
	AvailableBalance float64

	UserID uuid.UUID
	User   User

	PlaidItemID uuid.UUID
	PlaidItem   PlaidItem

	Transactions []Transaction `gorm:"foreignKey:BankAccountID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Transaction struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key"`
	Category    string    `json:"category"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	Type        string    `json:"type"` // E.g., "expense", "income"
	IsRecurring bool      `json:"is_recurring"`
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
