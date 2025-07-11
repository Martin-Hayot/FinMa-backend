package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"FinMa/internal/domain"
	"FinMa/internal/repository"
)

// BankAccountRepository implements the repository.BankAccountRepository interface
type BankAccountRepository struct {
	db *gorm.DB
}

// NewBankAccountRepository creates a new bank account repository
func NewBankAccountRepository(db *gorm.DB) *BankAccountRepository {
	return &BankAccountRepository{
		db: db,
	}
}

// Create adds a new bank account to the database
func (r *BankAccountRepository) Create(ctx context.Context, bankAccount *domain.BankAccount) error {
	return r.db.WithContext(ctx).Create(bankAccount).Error
}

// GetByID retrieves a bank account by ID with related data
func (r *BankAccountRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.BankAccount, error) {
	var bankAccount domain.BankAccount
	result := r.db.WithContext(ctx).
		Preload("User").
		Preload("Transactions").
		First(&bankAccount, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.BankAccount{}, errors.New("bank account not found")
		}
		return domain.BankAccount{}, result.Error
	}
	return bankAccount, nil
}

// GetByUserID retrieves all bank accounts for a specific user
func (r *BankAccountRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.BankAccount, error) {
	var bankAccounts []domain.BankAccount
	result := r.db.WithContext(ctx).
		Preload("Transactions").
		Where("user_id = ?", userID).
		Find(&bankAccounts)

	if result.Error != nil {
		return nil, result.Error
	}
	return bankAccounts, nil
}

// GetByAccountID retrieves a bank account by GoCardless account ID
func (r *BankAccountRepository) GetByAccountID(ctx context.Context, accountID string) (*domain.BankAccount, error) {
	var bankAccount domain.BankAccount
	result := r.db.WithContext(ctx).
		Preload("User").
		Preload("Requisition").
		Preload("Transactions").
		First(&bankAccount, "account_id = ?", accountID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, result.Error
	}
	return &bankAccount, nil
}



// Update updates a bank account in the database
func (r *BankAccountRepository) Update(ctx context.Context, bankAccount *domain.BankAccount) error {
	return r.db.WithContext(ctx).Save(bankAccount).Error
}

// Delete soft deletes a bank account from the database
func (r *BankAccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.BankAccount{}, "id = ?", id).Error
}

// ExistsByAccountID checks if a bank account with the given account ID exists
func (r *BankAccountRepository) ExistsByAccountID(ctx context.Context, accountID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.BankAccount{}).Where("account_id = ?", accountID).Count(&count).Error
	return count > 0, err
}



// GetUserAccountsWithBalance retrieves user's bank accounts with current balance information
func (r *BankAccountRepository) GetUserAccountsWithBalance(ctx context.Context, userID uuid.UUID) ([]domain.BankAccount, error) {
	var bankAccounts []domain.BankAccount
	result := r.db.WithContext(ctx).
		Select("id", "account_id", "name", "type", "currency", "institution_name",
			"balance_available", "balance_current", "iban", "user_id",
			"created_at", "updated_at").
		Where("user_id = ?", userID).
		Find(&bankAccounts)

	if result.Error != nil {
		return nil, result.Error
	}
	return bankAccounts, nil
}

// UpdateBalance updates the balance for a specific bank account
func (r *BankAccountRepository) UpdateBalance(ctx context.Context, accountID string, availableBalance, currentBalance float64) error {
	return r.db.WithContext(ctx).
		Model(&domain.BankAccount{}).
		Where("account_id = ?", accountID).
		Updates(map[string]interface{}{
			"balance_available": availableBalance,
			"balance_current":   currentBalance,
		}).Error
}

// CountByUserID returns the number of bank accounts for a user
func (r *BankAccountRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.BankAccount{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}


