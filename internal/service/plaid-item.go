package service

import (
	"FinMa/internal/domain"
	"FinMa/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type PlaidItemService interface {
	CreatePlaidItem(ctx context.Context, userID uuid.UUID, accessToken, itemID string) (domain.PlaidItem, error)
	CreatePlaidItemWithInstitution(ctx context.Context, userID uuid.UUID, accessToken, itemID, institution string) (domain.PlaidItem, error)
	GetPlaidItemByID(ctx context.Context, id uuid.UUID) (domain.PlaidItem, error)
	GetPlaidItemByItemID(ctx context.Context, itemID string) (domain.PlaidItem, error)
	GetPlaidItemByAccessToken(ctx context.Context, accessToken string) (domain.PlaidItem, error)
	GetUserPlaidItems(ctx context.Context, userID uuid.UUID) ([]domain.PlaidItem, error)
	UpdatePlaidItem(ctx context.Context, plaidItem *domain.PlaidItem) error
	DeletePlaidItem(ctx context.Context, id uuid.UUID) error
	ExistsByItemID(ctx context.Context, itemID string) (bool, error)
}

type plaidItemService struct {
	plaidItemRepo repository.PlaidItemRepository
	userRepo      repository.UserRepository
}

func NewPlaidItemService(
	plaidItemRepo repository.PlaidItemRepository,
	userRepo repository.UserRepository,
) PlaidItemService {
	return &plaidItemService{
		plaidItemRepo: plaidItemRepo,
		userRepo:      userRepo,
	}
}

// CreatePlaidItem creates a new Plaid item for a user
func (s *plaidItemService) CreatePlaidItem(ctx context.Context, userID uuid.UUID, accessToken, itemID string) (domain.PlaidItem, error) {
	return s.CreatePlaidItemWithInstitution(ctx, userID, accessToken, itemID, "")
}

// CreatePlaidItemWithInstitution creates a new Plaid item with institution information
func (s *plaidItemService) CreatePlaidItemWithInstitution(ctx context.Context, userID uuid.UUID, accessToken, itemID, institution string) (domain.PlaidItem, error) {
	// Verify user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return domain.PlaidItem{}, errors.New("user not found")
	}

	// Check if Plaid item already exists
	exists, err := s.plaidItemRepo.ExistsByItemID(ctx, itemID)
	if err != nil {
		return domain.PlaidItem{}, err
	}
	if exists {
		return domain.PlaidItem{}, errors.New("plaid item already exists")
	}

	// Create new Plaid item
	plaidItem := domain.PlaidItem{
		ID:          uuid.New(),
		ItemID:      itemID,
		AccessToken: accessToken,
		Institution: institution,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to database
	if err := s.plaidItemRepo.Create(ctx, &plaidItem); err != nil {
		return domain.PlaidItem{}, err
	}

	return plaidItem, nil
}

// GetPlaidItemByID retrieves a Plaid item by ID
func (s *plaidItemService) GetPlaidItemByID(ctx context.Context, id uuid.UUID) (domain.PlaidItem, error) {
	return s.plaidItemRepo.GetByID(ctx, id)
}

// GetPlaidItemByItemID retrieves a Plaid item by Plaid item ID
func (s *plaidItemService) GetPlaidItemByItemID(ctx context.Context, itemID string) (domain.PlaidItem, error) {
	return s.plaidItemRepo.GetByItemID(ctx, itemID)
}

// GetPlaidItemByAccessToken retrieves a Plaid item by access token
func (s *plaidItemService) GetPlaidItemByAccessToken(ctx context.Context, accessToken string) (domain.PlaidItem, error) {
	return s.plaidItemRepo.GetByAccessToken(ctx, accessToken)
}

// GetUserPlaidItems retrieves all Plaid items for a user
func (s *plaidItemService) GetUserPlaidItems(ctx context.Context, userID uuid.UUID) ([]domain.PlaidItem, error) {
	return s.plaidItemRepo.GetByUserID(ctx, userID)
}

// UpdatePlaidItem updates a Plaid item
func (s *plaidItemService) UpdatePlaidItem(ctx context.Context, plaidItem *domain.PlaidItem) error {
	plaidItem.UpdatedAt = time.Now()
	return s.plaidItemRepo.Update(ctx, plaidItem)
}

// DeletePlaidItem deletes a Plaid item
func (s *plaidItemService) DeletePlaidItem(ctx context.Context, id uuid.UUID) error {
	return s.plaidItemRepo.Delete(ctx, id)
}

// ExistsByItemID checks if a Plaid item exists by item ID
func (s *plaidItemService) ExistsByItemID(ctx context.Context, itemID string) (bool, error) {
	return s.plaidItemRepo.ExistsByItemID(ctx, itemID)
}
