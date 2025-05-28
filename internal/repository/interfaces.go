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
}

// BankAccountRepository defines operations for bank account data access
type BankAccountRepository interface {
	Create(ctx context.Context, bankAccount *domain.BankAccount) error
	Update(ctx context.Context, bankAccount *domain.BankAccount) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.BankAccount, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.BankAccount, error)
	GetByAccountID(ctx context.Context, accountID string) (domain.BankAccount, error)
	GetByGclItemID(ctx context.Context, gclItemID uuid.UUID) ([]domain.BankAccount, error)
	ExistsByAccountID(ctx context.Context, accountID string) (bool, error)
}

// GclItemRepository defines operations for GoCardless item data access
type GclItemRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, gclItem *domain.GclItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.GclItem, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.GclItem, error)
	GetByAccessToken(ctx context.Context, accessToken string) (*domain.GclItem, error)
	GetByProviderName(ctx context.Context, userID uuid.UUID, providerName string) (*domain.GclItem, error)
	Update(ctx context.Context, gclItem *domain.GclItem) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
