package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"FinMa/internal/domain"
)

type GoCardlessItemRepository struct {
	db *gorm.DB
}

// NewGoCardlessItemRepository creates a new GoCardless item repository
func NewGoCardlessItemRepository(db *gorm.DB) *GoCardlessItemRepository {
	return &GoCardlessItemRepository{
		db: db,
	}
}

// Create creates a new GoCardless item
func (r *GoCardlessItemRepository) Create(ctx context.Context, goCardlessItem *domain.GoCardlessItem) error {
	if err := r.db.WithContext(ctx).Create(goCardlessItem).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a GoCardless item by ID
func (r *GoCardlessItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.GoCardlessItem, error) {
	var goCardlessItem domain.GoCardlessItem
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&goCardlessItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &goCardlessItem, nil
}

// GetByUserID retrieves all GoCardless items for a user
func (r *GoCardlessItemRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.GoCardlessItem, error) {
	var goCardlessItems []domain.GoCardlessItem
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&goCardlessItems).Error; err != nil {
		return nil, err
	}
	return goCardlessItems, nil
}

// GetByAccessToken retrieves a GoCardless item by access token
func (r *GoCardlessItemRepository) GetByAccessToken(ctx context.Context, accessToken string) (*domain.GoCardlessItem, error) {
	var goCardlessItem domain.GoCardlessItem
	if err := r.db.WithContext(ctx).Where("access_token = ?", accessToken).First(&goCardlessItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &goCardlessItem, nil
}

// GetByProviderName retrieves a GoCardless item by provider name for a specific user
func (r *GoCardlessItemRepository) GetByProviderName(ctx context.Context, userID uuid.UUID, providerName string) (*domain.GoCardlessItem, error) {
	var goCardlessItem domain.GoCardlessItem
	if err := r.db.WithContext(ctx).Where("user_id = ? AND provider_name = ?", userID, providerName).First(&goCardlessItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &goCardlessItem, nil
}

// Update updates a GoCardless item
func (r *GoCardlessItemRepository) Update(ctx context.Context, goCardlessItem *domain.GoCardlessItem) error {
	if err := r.db.WithContext(ctx).Save(goCardlessItem).Error; err != nil {
		return err
	}
	return nil
}

// UpdateTokens updates the access and refresh tokens for a GoCardless item
func (r *GoCardlessItemRepository) UpdateTokens(ctx context.Context, id uuid.UUID, accessToken, refreshToken string, expiresAt *time.Time) error {
	updates := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"updated_at":    time.Now(),
	}

	if expiresAt != nil {
		updates["expires_at"] = *expiresAt
	}

	if err := r.db.WithContext(ctx).Model(&domain.GoCardlessItem{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

// Delete deletes a GoCardless item by ID
func (r *GoCardlessItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.GoCardlessItem{}, id).Error; err != nil {
		return err
	}
	return nil
}

// DeleteByUserID deletes all GoCardless items for a user
func (r *GoCardlessItemRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.GoCardlessItem{}).Error; err != nil {
		return err
	}
	return nil
}

// GetExpiredItems retrieves all GoCardless items with expired tokens
func (r *GoCardlessItemRepository) GetExpiredItems(ctx context.Context) ([]domain.GoCardlessItem, error) {
	var goCardlessItems []domain.GoCardlessItem
	if err := r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Find(&goCardlessItems).Error; err != nil {
		return nil, err
	}
	return goCardlessItems, nil
}

// GetItemsNeedingRefresh retrieves GoCardless items that need token refresh (expiring soon)
func (r *GoCardlessItemRepository) GetItemsNeedingRefresh(ctx context.Context, refreshThreshold time.Duration) ([]domain.GoCardlessItem, error) {
	var goCardlessItems []domain.GoCardlessItem
	refreshTime := time.Now().Add(refreshThreshold)

	if err := r.db.WithContext(ctx).Where("expires_at < ? AND expires_at > ?", refreshTime, time.Now()).Find(&goCardlessItems).Error; err != nil {
		return nil, err
	}
	return goCardlessItems, nil
}

// UpdateLastSyncTime updates the last sync time for a GoCardless item
func (r *GoCardlessItemRepository) UpdateLastSyncTime(ctx context.Context, id uuid.UUID, lastSyncTime time.Time) error {
	if err := r.db.WithContext(ctx).Model(&domain.GoCardlessItem{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_sync_time": lastSyncTime,
		"updated_at":     time.Now(),
	}).Error; err != nil {
		return err
	}
	return nil
}

// GetItemsForSync retrieves GoCardless items that need data synchronization
func (r *GoCardlessItemRepository) GetItemsForSync(ctx context.Context, syncInterval time.Duration) ([]domain.GoCardlessItem, error) {
	var goCardlessItems []domain.GoCardlessItem
	syncThreshold := time.Now().Add(-syncInterval)

	query := r.db.WithContext(ctx).Where("last_sync_time < ? OR last_sync_time IS NULL", syncThreshold)

	if err := query.Find(&goCardlessItems).Error; err != nil {
		return nil, err
	}
	return goCardlessItems, nil
}

// CountByUserID counts the number of GoCardless items for a user
func (r *GoCardlessItemRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.GoCardlessItem{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetWithBankAccounts retrieves a GoCardless item with its associated bank accounts
func (r *GoCardlessItemRepository) GetWithBankAccounts(ctx context.Context, id uuid.UUID) (*domain.GoCardlessItem, error) {
	var goCardlessItem domain.GoCardlessItem
	if err := r.db.WithContext(ctx).Preload("BankAccounts").Where("id = ?", id).First(&goCardlessItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &goCardlessItem, nil
}

// GetUserItemsWithBankAccounts retrieves all GoCardless items for a user with their bank accounts
func (r *GoCardlessItemRepository) GetUserItemsWithBankAccounts(ctx context.Context, userID uuid.UUID) ([]domain.GoCardlessItem, error) {
	var goCardlessItems []domain.GoCardlessItem
	if err := r.db.WithContext(ctx).Preload("BankAccounts").Where("user_id = ?", userID).Find(&goCardlessItems).Error; err != nil {
		return nil, err
	}
	return goCardlessItems, nil
}
