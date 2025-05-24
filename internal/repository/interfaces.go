package repository

import (
	"context"

	"github.com/google/uuid"

	"FinMa/internal/domain"
)

// UserRepository defines operations for user data access
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	SavePlaidAccessToken(ctx context.Context, userID uuid.UUID, accessToken, itemID string) error
}

// BankAccountRepository defines operations for bank account data access
type BankAccountRepository interface {
	Create(ctx context.Context, bankAccount *domain.BankAccount) error
	Update(ctx context.Context, bankAccount *domain.BankAccount) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.BankAccount, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.BankAccount, error)
	GetByAccountID(ctx context.Context, accountID string) (domain.BankAccount, error)
	GetByPlaidItemID(ctx context.Context, plaidItemID uuid.UUID) ([]domain.BankAccount, error)
	ExistsByAccountID(ctx context.Context, accountID string) (bool, error)
}

// PlaidItemRepository defines operations for Plaid item data access
type PlaidItemRepository interface {
	Create(ctx context.Context, plaidItem *domain.PlaidItem) error
	Update(ctx context.Context, plaidItem *domain.PlaidItem) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.PlaidItem, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.PlaidItem, error)
	GetByItemID(ctx context.Context, itemID string) (domain.PlaidItem, error)
	GetByAccessToken(ctx context.Context, accessToken string) (domain.PlaidItem, error)
	ExistsByItemID(ctx context.Context, itemID string) (bool, error)
}
