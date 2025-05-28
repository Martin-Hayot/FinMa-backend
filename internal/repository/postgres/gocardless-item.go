package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"FinMa/internal/domain"
)

type GclItemRepository struct {
	db *gorm.DB
}

// NewGclItemRepository creates a new Gcl item repository
func NewGclItemRepository(db *gorm.DB) *GclItemRepository {
	return &GclItemRepository{
		db: db,
	}
}

// Create creates a new GoCardless item
func (r *GclItemRepository) Create(ctx context.Context, gclItem *domain.GclItem) error {
	if err := r.db.WithContext(ctx).Create(gclItem).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a Gcl item by ID
func (r *GclItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.GclItem, error) {
	var gclItem domain.GclItem
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&gclItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &gclItem, nil
}

// GetByUserID retrieves all Gcl items for a user
func (r *GclItemRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.GclItem, error) {
	var gclItems []domain.GclItem
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&gclItems).Error; err != nil {
		return nil, err
	}
	return gclItems, nil
}

// GetByAccessToken retrieves a Gcl item by access token
func (r *GclItemRepository) GetByAccessToken(ctx context.Context, accessToken string) (*domain.GclItem, error) {
	var gclItem domain.GclItem
	if err := r.db.WithContext(ctx).Where("access_token = ?", accessToken).First(&gclItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &gclItem, nil
}

// GetByProviderName retrieves a Gcl item by provider name for a specific user
func (r *GclItemRepository) GetByProviderName(ctx context.Context, userID uuid.UUID, providerName string) (*domain.GclItem, error) {
	var gclItem domain.GclItem
	if err := r.db.WithContext(ctx).Where("user_id = ? AND provider_name = ?", userID, providerName).First(&gclItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &gclItem, nil
}

// Update updates a Gcl item
func (r *GclItemRepository) Update(ctx context.Context, gclItem *domain.GclItem) error {
	if err := r.db.WithContext(ctx).Save(gclItem).Error; err != nil {
		return err
	}
	return nil
}

// UpdateTokens updates the access and refresh tokens for a Gcl item
func (r *GclItemRepository) UpdateTokens(ctx context.Context, id uuid.UUID, accessToken, refreshToken string, expiresAt *time.Time) error {
	updates := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"updated_at":    time.Now(),
	}

	if expiresAt != nil {
		updates["expires_at"] = *expiresAt
	}

	if err := r.db.WithContext(ctx).Model(&domain.GclItem{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

// Delete deletes a Gcl item by ID
func (r *GclItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.GclItem{}, id).Error; err != nil {
		return err
	}
	return nil
}

// DeleteByUserID deletes all Gcl items for a user
func (r *GclItemRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.GclItem{}).Error; err != nil {
		return err
	}
	return nil
}
