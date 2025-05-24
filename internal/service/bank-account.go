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
	CreateBankAccount(ctx context.Context, userID uuid.UUID, plaidItemID uuid.UUID, accountData domain.BankAccount) (domain.BankAccount, error)
	CreateBankAccountsFromPlaid(ctx context.Context, userID uuid.UUID, plaidItemID uuid.UUID, accounts []domain.BankAccount) ([]domain.BankAccount, error)
	GetUserBankAccounts(ctx context.Context, userID uuid.UUID) ([]domain.BankAccount, error)
	GetBankAccountByID(ctx context.Context, id uuid.UUID) (domain.BankAccount, error)
	GetBankAccountsByPlaidItem(ctx context.Context, plaidItemID uuid.UUID) ([]domain.BankAccount, error)
	UpdateBankAccount(ctx context.Context, bankAccount *domain.BankAccount) error
	DeleteBankAccount(ctx context.Context, id uuid.UUID) error
	SyncBankAccountBalances(ctx context.Context, plaidItemID uuid.UUID, accountBalances map[string]domain.BankAccount) error
}

type bankAccountService struct {
	bankAccountRepo repository.BankAccountRepository
	plaidItemRepo   repository.PlaidItemRepository
	userRepo        repository.UserRepository
}

// NewBankAccountService creates a new bank account service
func NewBankAccountService(
	bankAccountRepo repository.BankAccountRepository,
	plaidItemRepo repository.PlaidItemRepository,
	userRepo repository.UserRepository,
) BankAccountService {
	return &bankAccountService{
		bankAccountRepo: bankAccountRepo,
		plaidItemRepo:   plaidItemRepo,
		userRepo:        userRepo,
	}
}

// CreateBankAccount creates a new bank account for a user
func (s *bankAccountService) CreateBankAccount(ctx context.Context, userID uuid.UUID, plaidItemID uuid.UUID, accountData domain.BankAccount) (domain.BankAccount, error) {
	// Verify user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return domain.BankAccount{}, errors.New("user not found")
	}

	// Verify Plaid item exists and belongs to user
	plaidItem, err := s.plaidItemRepo.GetByID(ctx, plaidItemID)
	if err != nil {
		return domain.BankAccount{}, errors.New("plaid item not found")
	}
	if plaidItem.UserID != userID {
		return domain.BankAccount{}, errors.New("plaid item does not belong to user")
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
	accountData.PlaidItemID = plaidItemID
	accountData.CreatedAt = time.Now()
	accountData.UpdatedAt = time.Now()

	// Create the bank account
	if err := s.bankAccountRepo.Create(ctx, &accountData); err != nil {
		return domain.BankAccount{}, err
	}

	return accountData, nil
}

// CreateBankAccountsFromPlaid creates multiple bank accounts from Plaid data
func (s *bankAccountService) CreateBankAccountsFromPlaid(ctx context.Context, userID uuid.UUID, plaidItemID uuid.UUID, accounts []domain.BankAccount) ([]domain.BankAccount, error) {
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
		createdAccount, err := s.CreateBankAccount(ctx, userID, plaidItemID, account)
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

// GetBankAccountsByPlaidItem retrieves all bank accounts for a Plaid item
func (s *bankAccountService) GetBankAccountsByPlaidItem(ctx context.Context, plaidItemID uuid.UUID) ([]domain.BankAccount, error) {
	return s.bankAccountRepo.GetByPlaidItemID(ctx, plaidItemID)
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

// SyncBankAccountBalances updates account balances from Plaid data
func (s *bankAccountService) SyncBankAccountBalances(ctx context.Context, plaidItemID uuid.UUID, accountBalances map[string]domain.BankAccount) error {
	// Get existing accounts for this Plaid item
	existingAccounts, err := s.bankAccountRepo.GetByPlaidItemID(ctx, plaidItemID)
	if err != nil {
		return err
	}

	// Update balances for existing accounts
	for _, account := range existingAccounts {
		if updatedData, exists := accountBalances[account.AccountID]; exists {
			account.CurrentBalance = updatedData.CurrentBalance
			account.AvailableBalance = updatedData.AvailableBalance
			account.UpdatedAt = time.Now()

			if err := s.bankAccountRepo.Update(ctx, &account); err != nil {
				continue // Continue with other accounts if one fails
			}
		}
	}

	return nil
}
