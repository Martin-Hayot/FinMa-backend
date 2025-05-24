package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"FinMa/internal/domain"
)

// PlaidItemRepository implements the repository.PlaidItemRepository interface
type PlaidItemRepository struct {
	db *gorm.DB
}

// NewPlaidItemRepository creates a new Plaid item repository
func NewPlaidItemRepository(db *gorm.DB) *PlaidItemRepository {
	return &PlaidItemRepository{
		db: db,
	}
}

// Create adds a new Plaid item to the database
func (r *PlaidItemRepository) Create(ctx context.Context, plaidItem *domain.PlaidItem) error {
	return r.db.WithContext(ctx).Create(plaidItem).Error
}

// GetByID retrieves a Plaid item by ID
func (r *PlaidItemRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.PlaidItem, error) {
	var plaidItem domain.PlaidItem
	result := r.db.WithContext(ctx).
		Preload("User").
		First(&plaidItem, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.PlaidItem{}, errors.New("plaid item not found")
		}
		return domain.PlaidItem{}, result.Error
	}
	return plaidItem, nil
}

// GetByUserID retrieves all Plaid items for a specific user
func (r *PlaidItemRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.PlaidItem, error) {
	var plaidItems []domain.PlaidItem
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&plaidItems)

	if result.Error != nil {
		return nil, result.Error
	}
	return plaidItems, nil
}

// GetByItemID retrieves a Plaid item by Plaid item ID
func (r *PlaidItemRepository) GetByItemID(ctx context.Context, itemID string) (domain.PlaidItem, error) {
	var plaidItem domain.PlaidItem
	result := r.db.WithContext(ctx).
		Preload("User").
		First(&plaidItem, "item_id = ?", itemID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.PlaidItem{}, errors.New("plaid item not found")
		}
		return domain.PlaidItem{}, result.Error
	}
	return plaidItem, nil
}

// GetByAccessToken retrieves a Plaid item by access token
func (r *PlaidItemRepository) GetByAccessToken(ctx context.Context, accessToken string) (domain.PlaidItem, error) {
	var plaidItem domain.PlaidItem
	result := r.db.WithContext(ctx).
		Preload("User").
		First(&plaidItem, "access_token = ?", accessToken)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.PlaidItem{}, errors.New("plaid item not found")
		}
		return domain.PlaidItem{}, result.Error
	}
	return plaidItem, nil
}

// Update updates a Plaid item in the database
func (r *PlaidItemRepository) Update(ctx context.Context, plaidItem *domain.PlaidItem) error {
	return r.db.WithContext(ctx).Save(plaidItem).Error
}

// Delete deletes a Plaid item from the database
func (r *PlaidItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.PlaidItem{}, "id = ?", id).Error
}

// ExistsByItemID checks if a Plaid item with the given item ID exists
func (r *PlaidItemRepository) ExistsByItemID(ctx context.Context, itemID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.PlaidItem{}).Where("item_id = ?", itemID).Count(&count).Error
	return count > 0, err
}
