package database

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uint           `json:"id" gorm:"primary_key"`
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	Email        string         `json:"email" gorm:"uniqueIndex"`
	Password     string         `json:"password"`
	Role         string         `json:"role"`
	Transactions []Transactions `json:"transactions"`
	BankAccounts []BankAccount  `json:"bank_accounts"`
}

type BankAccount struct {
	gorm.Model
	ID            uint    `json:"id" gorm:"primary_key"`
	BankName      string  `json:"bank_name"`
	AccountType   string  `json:"account_type"`
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`

	UserID       uint           `json:"user_id"`
	User         User           `json:"user"`
	Transactions []Transactions `json:"transactions"`
}
type Transactions struct {
	gorm.Model
	ID          uint      `json:"id" gorm:"primary_key"`
	Category    string    `json:"category"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`
	IsRecuring  bool      `json:"is_recuring"`
	Description string    `json:"description"`

	UserID        uint        `json:"user_id"`
	User          User        `json:"user"`
	BankAccountID uint        `json:"bank_account_id"`
	BankAccount   BankAccount `json:"bank_account"`
}

type Budget struct {
	gorm.Model
	ID        uint      `json:"id" gorm:"primary_key"`
	Category  string    `json:"category"`
	Amount    float64   `json:"amount"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`

	UserID uint `json:"user_id"`
	User   User `json:"user"`
}

type Notification struct {
	gorm.Model
	ID       uint   `json:"id" gorm:"primary_key"`
	Type     string `json:"type"`
	Message  string `json:"message"`
	IsActive bool   `json:"is_active"`

	UserID uint `json:"user_id"`
	User   User `json:"user"`
}
