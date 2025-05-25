package repository

import (
	"context"
	"time"

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
	GetByGoCardlessItemID(ctx context.Context, goCardlessItemID uuid.UUID) ([]domain.BankAccount, error)
	ExistsByAccountID(ctx context.Context, accountID string) (bool, error)
}

// GoCardlessItemRepository defines operations for GoCardless item data access
type GoCardlessItemRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, goCardlessItem *domain.GoCardlessItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.GoCardlessItem, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.GoCardlessItem, error)
	GetByAccessToken(ctx context.Context, accessToken string) (*domain.GoCardlessItem, error)
	GetByProviderName(ctx context.Context, userID uuid.UUID, providerName string) (*domain.GoCardlessItem, error)
	Update(ctx context.Context, goCardlessItem *domain.GoCardlessItem) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error

	// Token management
	UpdateTokens(ctx context.Context, id uuid.UUID, accessToken, refreshToken string, expiresAt *time.Time) error
	GetExpiredItems(ctx context.Context) ([]domain.GoCardlessItem, error)
	GetItemsNeedingRefresh(ctx context.Context, refreshThreshold time.Duration) ([]domain.GoCardlessItem, error)

	// Sync management
	UpdateLastSyncTime(ctx context.Context, id uuid.UUID, lastSyncTime time.Time) error
	GetItemsForSync(ctx context.Context, syncInterval time.Duration) ([]domain.GoCardlessItem, error)

	// Utility methods
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	GetWithBankAccounts(ctx context.Context, id uuid.UUID) (*domain.GoCardlessItem, error)
	GetUserItemsWithBankAccounts(ctx context.Context, userID uuid.UUID) ([]domain.GoCardlessItem, error)
}
