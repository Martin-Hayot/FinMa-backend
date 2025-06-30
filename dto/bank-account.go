package dto

import (
	"time"

	"github.com/google/uuid"
)

// BankAccountResponse represents the data returned for a bank account to the client
type BankAccountResponse struct {
	ID               uuid.UUID `json:"id"`
	AccountID        string    `json:"account_id"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	Currency         string    `json:"currency"`
	InstitutionName  string    `json:"institution_name"`
	BalanceAvailable float64   `json:"balance_available"`
	BalanceCurrent   float64   `json:"balance_current"`
	IBAN             string    `json:"iban,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}