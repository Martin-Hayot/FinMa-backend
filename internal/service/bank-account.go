package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"FinMa/internal/domain"
	"FinMa/internal/repository"
)

// BankAccountService defines operations for bank account management
type BankAccountService interface {
	CreateBankAccount(ctx context.Context, userID uuid.UUID, goCardlessItemID uuid.UUID, accountData domain.BankAccount) (domain.BankAccount, error)
	GetUserBankAccounts(ctx context.Context, userID uuid.UUID) ([]domain.BankAccount, error)
	GetBankAccountByID(ctx context.Context, id uuid.UUID) (domain.BankAccount, error)
	GetBankAccountsByGoCardlessItem(ctx context.Context, goCardlessItemID uuid.UUID) ([]domain.BankAccount, error)
	UpdateBankAccount(ctx context.Context, bankAccount *domain.BankAccount) error
	DeleteBankAccount(ctx context.Context, id uuid.UUID) error
}

type bankAccountService struct {
	bankAccountRepo    repository.BankAccountRepository
	goCardlessItemRepo repository.GoCardlessItemRepository
	userRepo           repository.UserRepository
}

// NewBankAccountService creates a new bank account service
func NewBankAccountService(
	bankAccountRepo repository.BankAccountRepository,
	goCardlessItemRepo repository.GoCardlessItemRepository,
	userRepo repository.UserRepository,
) BankAccountService {
	return &bankAccountService{
		bankAccountRepo:    bankAccountRepo,
		goCardlessItemRepo: goCardlessItemRepo,
		userRepo:           userRepo,
	}
}

// CreateBankAccount creates a new bank account for a user
func (s *bankAccountService) CreateBankAccount(ctx context.Context, userID uuid.UUID, goCardlessItemID uuid.UUID, accountData domain.BankAccount) (domain.BankAccount, error) {
	// Verify user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return domain.BankAccount{}, errors.New("user not found")
	}

	// Verify GoCardless item exists and belongs to user
	goCardlessItem, err := s.goCardlessItemRepo.GetByID(ctx, goCardlessItemID)
	if err != nil {
		return domain.BankAccount{}, errors.New("gocardless item not found")
	}
	if goCardlessItem.UserID != userID {
		return domain.BankAccount{}, errors.New("gocardless item does not belong to user")
	}

	// Check if account already exists
	exists, err := s.bankAccountRepo.ExistsByAccountID(ctx, accountData.AccountID)
	if err != nil {
		return domain.BankAccount{}, err
	}
	if exists {
		return domain.BankAccount{}, errors.New("bank account already exists")
	}

	// Set required fields
	accountData.ID = uuid.New()
	accountData.UserID = userID
	accountData.GoCardlessItemID = goCardlessItemID
	accountData.CreatedAt = time.Now()
	accountData.UpdatedAt = time.Now()

	// Create the bank account
	if err := s.bankAccountRepo.Create(ctx, &accountData); err != nil {
		return domain.BankAccount{}, err
	}

	return accountData, nil
}

// CreateBankAccountsFromGoCardless creates multiple bank accounts from GoCardless data
func (s *bankAccountService) CreateBankAccountsFromGoCardless(ctx context.Context, userID uuid.UUID, goCardlessItemID uuid.UUID, accounts []domain.BankAccount) ([]domain.BankAccount, error) {
	var createdAccounts []domain.BankAccount

	for _, account := range accounts {
		// Check if account already exists
		exists, err := s.bankAccountRepo.ExistsByAccountID(ctx, account.AccountID)
		if err != nil {
			continue // Skip this account and continue with others
		}
		if exists {
			continue // Skip existing accounts
		}

		// Create the account
		createdAccount, err := s.CreateBankAccount(ctx, userID, goCardlessItemID, account)
		if err != nil {
			continue // Skip failed accounts and continue with others
		}

		createdAccounts = append(createdAccounts, createdAccount)
	}

	return createdAccounts, nil
}

// GetUserBankAccounts retrieves all bank accounts for a user
func (s *bankAccountService) GetUserBankAccounts(ctx context.Context, userID uuid.UUID) ([]domain.BankAccount, error) {
	return s.bankAccountRepo.GetByUserID(ctx, userID)
}

// GetBankAccountByID retrieves a bank account by ID
func (s *bankAccountService) GetBankAccountByID(ctx context.Context, id uuid.UUID) (domain.BankAccount, error) {
	return s.bankAccountRepo.GetByID(ctx, id)
}

// GetBankAccountsByGoCardlessItem retrieves all bank accounts for a gocardless item
func (s *bankAccountService) GetBankAccountsByGoCardlessItem(ctx context.Context, goCardlessItemID uuid.UUID) ([]domain.BankAccount, error) {
	return s.bankAccountRepo.GetByGoCardlessItemID(ctx, goCardlessItemID)
}

// UpdateBankAccount updates a bank account
func (s *bankAccountService) UpdateBankAccount(ctx context.Context, bankAccount *domain.BankAccount) error {
	bankAccount.UpdatedAt = time.Now()
	return s.bankAccountRepo.Update(ctx, bankAccount)
}

// DeleteBankAccount deletes a bank account
func (s *bankAccountService) DeleteBankAccount(ctx context.Context, id uuid.UUID) error {
	return s.bankAccountRepo.Delete(ctx, id)
}
