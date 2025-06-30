package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	Requisitions  []Requisition  `gorm:"foreignKey:UserID"`
	BankAccounts  []BankAccount  `gorm:"foreignKey:UserID"`
	Transactions  []Transaction  `gorm:"foreignKey:UserID"`
	Budgets       []Budget       `gorm:"foreignKey:UserID"`
	Notifications []Notification `gorm:"foreignKey:UserID"`
	RefreshTokens []RefreshToken `gorm:"foreignKey:UserID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete
}

type Requisition struct {
	ID            string `gorm:"primaryKey;not null" json:"id"` // GoCardless Requisition ID as primary key
	Status        string `gorm:"not null" json:"status"`        // e.g., "created", "redirected", "linked", "expired"
	RedirectURI   string `gorm:"not null" json:"redirect_uri"`
	InstitutionID string `gorm:"not null" json:"institution_id"`
	Link          string `json:"link"`      // Authorization link provided by GoCardless
	Reference     string `json:"reference"` // Unique ID for internal reference

	UserID uuid.UUID `gorm:"not null;index" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID" json:"-"`

	// Association with bank accounts
	BankAccounts []BankAccount `gorm:"foreignKey:RequisitionID" json:"bank_accounts,omitempty"`

	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
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
	ID               uuid.UUID `gorm:"primaryKey" json:"id"`
	AccountID        string    `gorm:"uniqueIndex;not null" json:"account_id"` // GoCardless account ID
	Name             string    `gorm:"not null" json:"name"`
	Type             string    `gorm:"not null" json:"type"`
	Currency         string    `gorm:"not null" json:"currency"`
	InstitutionName  string    `json:"institution_name"`
	BalanceAvailable float64   `json:"balance_available"`
	BalanceCurrent   float64   `json:"balance_current"`
	IBAN             string    `json:"iban,omitempty"`

	UserID        uuid.UUID   `gorm:"not null;index" json:"user_id"`
	User          User        `gorm:"foreignKey:UserID" json:"-"`
	RequisitionID string      `gorm:"not null;index" json:"requisition_id"` // Link to requisition
	Requisition   Requisition `gorm:"foreignKey:RequisitionID" json:"-"`

	// Association with transactions
	Transactions []Transaction `gorm:"foreignKey:BankAccountID" json:"transactions,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
}
