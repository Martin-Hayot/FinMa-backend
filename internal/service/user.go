package service

import (
	"context"
	"errors"
	"time"

	"FinMa/dto"
	"FinMa/internal/domain"
	"FinMa/internal/repository"
	"FinMa/utils"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

// user related business logic
// UserService defines operations for user management
type UserService interface {
	// Profile management
	GetUserByID(ctx context.Context, id uuid.UUID) (dto.UserResponse, error)
	UpdateProfile(ctx context.Context, user domain.User, req dto.UpdateProfileRequest) (dto.UserResponse, error)
	// ChangePassword(ctx context.Context, id uuid.UUID, req dto.ChangePasswordRequest) error
	DeleteAccount(ctx context.Context, id uuid.UUID, password string) error

	// Plaid integration
	SavePlaidAccessToken(ctx context.Context, userID uuid.UUID, accessToken string, itemID string) error
	// Preference management
	// GetUserPreferences(ctx context.Context, userID uuid.UUID) (dto.UserPreferencesResponse, error)
	// UpdateUserPreferences(ctx context.Context, userID uuid.UUID, req dto.UpdatePreferencesRequest) error

	// User analytics
	// GetFinancialSummary(ctx context.Context, userID uuid.UUID) (dto.FinancialSummaryResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
	// preferenceRepo  repository.UserPreferenceRepository
	// transactionRepo repository.TransactionRepository
}

// UpdateProfile updates a user's profile information
func (s *userService) UpdateProfile(ctx context.Context, user domain.User, req dto.UpdateProfileRequest) (dto.UserResponse, error) {
	// Update user fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}

	// Check if email is being changed
	if req.Email != "" && req.Email != user.Email {
		// Check if the new email is already in use
		exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
		if err != nil {
			return dto.UserResponse{}, err
		}
		if exists {
			return dto.UserResponse{}, errors.New("email already in use")
		}
		user.Email = req.Email
	}

	user.UpdatedAt = time.Now()

	// Update user in database
	if err := s.userRepo.Update(ctx, &user); err != nil {
		log.Error("Failed to update user", "id", user.ID, "error", err)
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// DeleteAccount deletes a user account
func (s *userService) DeleteAccount(ctx context.Context, id uuid.UUID, password string) error {
	// Get the current user
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify the password
	if err := utils.ComparePasswords(user.Password, password); err != nil {
		return errors.New("incorrect password")
	}

	// Delete the user
	return s.userRepo.Delete(ctx, id)
}

// GetUserByID retrieves a user by their ID
func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		log.Error("Failed to get user", "id", id, "error", err)
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// SavePlaidAccessToken saves the Plaid access token and item ID for a user
func (s *userService) SavePlaidAccessToken(ctx context.Context, userID uuid.UUID, accessToken string, itemID string) error {
	// Validate access token and item ID
	if accessToken == "" || itemID == "" {
		return errors.New("invalid access token or item ID")
	}

	// Save the access token and item ID to the user repository
	return s.userRepo.SavePlaidAccessToken(ctx, userID, accessToken, itemID)
}

// NewUserService creates a new user service
func NewUserService(
	userRepo repository.UserRepository,
	// preferenceRepo repository.UserPreferenceRepository,
	// transactionRepo repository.TransactionRepository,
) UserService {
	return &userService{
		userRepo: userRepo,
		// preferenceRepo:  preferenceRepo,
		// transactionRepo: transactionRepo,
	}
}
