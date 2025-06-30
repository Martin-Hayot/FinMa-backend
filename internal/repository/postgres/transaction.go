package postgres

import (
	"context"
	"fmt"

	"FinMa/internal/domain"
	"FinMa/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) repository.TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	result := r.db.WithContext(ctx).Create(transaction)
	if result.Error != nil {
		return fmt.Errorf("failed to create transaction: %w", result.Error)
	}
	return nil
}

func (r *transactionRepository) CreateInBatches(ctx context.Context, transactions []*domain.Transaction) error {
	result := r.db.WithContext(ctx).CreateInBatches(transactions, 100) // Batch size of 100
	if result.Error != nil {
		return fmt.Errorf("failed to create transactions in batches: %w", result.Error)
	}
	return nil
}

func (r *transactionRepository) GetByBankAccountID(ctx context.Context, bankAccountID uuid.UUID) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	result := r.db.WithContext(ctx).Where("bank_account_id = ?", bankAccountID).Find(&transactions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get transactions by bank account ID: %w", result.Error)
	}
	return transactions, nil
}

func (r *transactionRepository) GetByTransactionID(ctx context.Context, transactionID string) (domain.Transaction, error) {
	var transaction domain.Transaction
	result := r.db.WithContext(ctx).Where("id = ?", transactionID).First(&transaction)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return domain.Transaction{}, repository.ErrNotFound
		}
		return domain.Transaction{}, fmt.Errorf("failed to get transaction by ID: %w", result.Error)
	}
	return transaction, nil
}

func (r *transactionRepository) ExistsByTransactionID(ctx context.Context, transactionID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Transaction{}).Where("id = ?", transactionID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check if transaction exists: %w", err)
	}
	return count > 0, nil
}
