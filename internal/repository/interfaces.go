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
	GetUserAccountsWithBalance(ctx context.Context, userID uuid.UUID) ([]domain.BankAccount, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.BankAccount, error)
	GetByAccountID(ctx context.Context, accountID string) (*domain.BankAccount, error)
	ExistsByAccountID(ctx context.Context, accountID string) (bool, error)
}

type RequisitionRepository interface {
	// Create adds a new requisition to the database
	Create(ctx context.Context, requisition *domain.Requisition) error
	// Update modifies an existing requisition
	Update(ctx context.Context, requisition *domain.Requisition) error
	// GetByID retrieves a requisition by its ID
	GetByID(ctx context.Context, id string) (*domain.Requisition, error)
	// GetByReference retrieves a requisition by its reference
	GetByReference(ctx context.Context, reference string) (*domain.Requisition, error)
	// GetByUserID retrieves all requisitions for a user
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Requisition, error)
	// GetByUserIDAndInstitutionID retrieves a requisition by user ID and institution ID
	GetByUserIDAndInstitutionID(ctx context.Context, userID uuid.UUID, institutionID string) (*domain.Requisition, error)
}

// TransactionRepository defines operations for transaction data access
type TransactionRepository interface {
	Create(ctx context.Context, transaction *domain.Transaction) error
	CreateInBatches(ctx context.Context, transactions []*domain.Transaction) error
	GetByBankAccountID(ctx context.Context, bankAccountID uuid.UUID) ([]domain.Transaction, error)
	GetByTransactionID(ctx context.Context, transactionID string) (domain.Transaction, error)
	ExistsByTransactionID(ctx context.Context, transactionID string) (bool, error)
}
