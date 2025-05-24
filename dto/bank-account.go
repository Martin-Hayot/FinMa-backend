package dto

import "github.com/google/uuid"

type BankAccountResponse struct {
	ID           uuid.UUID `json:"id"`
	AccountID    string    `json:"account_id"`
	AccountName  string    `json:"account_name"`
	OfficialName string    `json:"official_name"`
	AccountType  string    `json:"account_type"`
	Mask         string    `json:"mask"`
	PlaidItemID  uuid.UUID `json:"plaid_item_id"`
}
