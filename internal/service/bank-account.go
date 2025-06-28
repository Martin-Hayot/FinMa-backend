package service

import (
	"context"

	"github.com/google/uuid"

	"FinMa/internal/domain"
	"FinMa/internal/repository"
)

// BankAccountService defines operations for bank account management
type BankAccountService interface {
	CreateBankAccount(ctx context.Context, userID uuid.UUID, goCardlessItemID uuid.UUID, accountData domain.BankAccount) (domain.BankAccount, error)
	GetUserBankAccounts(ctx context.Context, userID uuid.UUID) ([]domain.BankAccount, error)
	GetBankAccountByID(ctx context.Context, id uuid.UUID) (domain.BankAccount, error)
	GetBankAccountsByGclItem(ctx context.Context, gclItemID uuid.UUID) ([]domain.BankAccount, error)
	UpdateBankAccount(ctx context.Context, bankAccount *domain.BankAccount) error
	DeleteBankAccount(ctx context.Context, id uuid.UUID) error
}

type bankAccountService struct {
	bankAccountRepo repository.BankAccountRepository
	userRepo        repository.UserRepository
}

// NewBankAccountService creates a new bank account service
// func NewBankAccountService(
// 	bankAccountRepo repository.BankAccountRepository,
// 	userRepo repository.UserRepository,
// ) BankAccountService {
// 	return &bankAccountService{
// 		bankAccountRepo: bankAccountRepo,
// 		userRepo:        userRepo,
// 	}
// }
