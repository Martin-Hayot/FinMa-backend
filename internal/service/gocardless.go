package service

import (
	"context"
	"fmt"
	"time"

	"FinMa/dto"
	"FinMa/internal/domain"
	"FinMa/internal/repository"
	"FinMa/pkg/gocardless"

	"github.com/google/uuid"
)

type GclService interface {
	// Token management
	GetValidAccessToken() (string, error)
	RefreshTokenIfNeeded() error
	GetTokenStatus() map[string]interface{}
	ClearToken()

	// LinkAccount initiates the linking of a bank account for a user with a specific institution
	LinkAccount(ctx context.Context, userID uuid.UUID, institutionID, redirectURL string) (*dto.LinkAccountResponse, error)

	// SyncRequisition syncs an existing requisition for a user
	SyncRequisition(ctx context.Context, requisitionReference string, userID uuid.UUID) (*dto.GoCardlessUpdateRequisitionResponse, error)

	// Institutions
	GetInstitutions(ctx context.Context, countryCode string) ([]dto.Institution, error)
}

type gclService struct {
	bankAccountRepo repository.BankAccountRepository
	userRepo        repository.UserRepository
	requisitionRepo repository.RequisitionRepository
	gclClient       *gocardless.Client
}

// NewGclService creates a new GoCardless service
func NewGclService(
	bankAccountRepo repository.BankAccountRepository,
	userRepo repository.UserRepository,
	requisitionRepo repository.RequisitionRepository,
	gclClient *gocardless.Client,
) GclService {
	return &gclService{
		bankAccountRepo: bankAccountRepo,
		requisitionRepo: requisitionRepo,
		userRepo:        userRepo,
		gclClient:       gclClient,
	}
}

func (s *gclService) LinkAccount(ctx context.Context, userID uuid.UUID, institutionID, redirectURL string) (*dto.LinkAccountResponse, error) {
	// verify that there is not a already an active requisition for this specific user and institution
	existingRequisition, err := s.requisitionRepo.GetByUserIDAndInstitutionID(ctx, userID, institutionID)
	if err != nil {
		// Only return error if it's not a "not found" error
		if !repository.IsNotFoundError(err) {
			return nil, fmt.Errorf("failed to check existing requisition: %w", err)
		}
		// If it's a "not found" error, continue with existingRequisition = nil
		existingRequisition = nil
	}
	// If there is an existing requisition, return the link
	if existingRequisition != nil {
		return &dto.LinkAccountResponse{
			Link: existingRequisition.Link,
		}, nil
	}

	// Call the GoCardless client to create a requisition
	response, err := s.gclClient.CreateRequisition(userID, institutionID, redirectURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create requisition: %w", err)
	}

	// store the requisition in the database
	err = s.requisitionRepo.Create(ctx, &domain.Requisition{
		ID:            response.ID,
		UserID:        userID,
		InstitutionID: institutionID,
		RedirectURI:   redirectURL,
		Status:        response.Status,
		Link:          response.Link,
		Reference:     response.Reference,
		ExpiresAt:     func(t time.Time) *time.Time { return &t }(time.Now().Add(90 * 24 * time.Hour)), // Set expiry to 90 days from now
	})
	if err != nil {
		return nil, fmt.Errorf("failed to store requisition: %w", err)
	}

	// Return the link to redirect the user to for linking their account
	return &dto.LinkAccountResponse{
		Link: response.Link,
	}, nil
}

func (s *gclService) SyncRequisition(ctx context.Context, requisitionReference string, userID uuid.UUID) (*dto.GoCardlessUpdateRequisitionResponse, error) {
	// Get the requisition by reference
	requisition, err := s.requisitionRepo.GetByReference(ctx, requisitionReference)
	if err != nil {
		return nil, fmt.Errorf("failed to get requisition by reference: %w", err)
	}

	// Call the GoCardless client to update the requisition
	response, err := s.gclClient.GetRequisition(userID, requisition.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update requisition: %w", err)
	}

	// Update the requisition in the database
	expiresAt := time.Now().Add(90 * 24 * time.Hour)
	err = s.requisitionRepo.Update(ctx, &domain.Requisition{
		ID:            response.ID,
		UserID:        requisition.UserID,
		InstitutionID: requisition.InstitutionID,
		RedirectURI:   requisition.RedirectURI,
		Status:        response.Status,
		Link:          response.Link,
		Reference:     response.Reference,
		ExpiresAt:     &expiresAt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update requisition in database: %w", err)
	}

	// Process account IDs if they exist in the response
	if len(response.Accounts) > 0 {
		err = s.processAccountsFromRequisition(ctx, response.Accounts, requisition.ID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to process accounts: %w", err)
		}
	}

	return &dto.GoCardlessUpdateRequisitionResponse{
		Status:        response.Status,
		InstitutionID: response.InstitutionID,
		Reference:     response.Reference,
	}, nil
}

// processAccountsFromRequisition processes the account IDs from a requisition
func (s *gclService) processAccountsFromRequisition(ctx context.Context, accountIDs []string, requisitionID string, userID uuid.UUID) error {
	for _, accountID := range accountIDs {
		// Check if the account already exists in our database
		existingAccount, err := s.bankAccountRepo.GetByAccountID(ctx, accountID)
		if err != nil && !repository.IsNotFoundError(err) {
			return fmt.Errorf("failed to check existing account: %w", err)
		}

		// If account doesn't exist, fetch details from GoCardless and create it
		if existingAccount == nil {
			accountDetails, err := s.gclClient.GetAccountDetails(accountID)
			if err != nil {
				return fmt.Errorf("failed to get account details for %s: %w", accountID, err)
			}

			// Create new bank account record
			bankAccount := &domain.BankAccount{
				ID:               uuid.New(),
				AccountID:        accountID,
				Name:             accountDetails.Name,
				Type:             accountDetails.Type,
				Currency:         accountDetails.Currency,
				InstitutionName:  accountDetails.InstitutionName,
				BalanceAvailable: accountDetails.BalanceAvailable,
				BalanceCurrent:   accountDetails.BalanceCurrent,
				IBAN:             accountDetails.IBAN,
				UserID:           userID,
				RequisitionID:    requisitionID,
			}

			err = s.bankAccountRepo.Create(ctx, bankAccount)
			if err != nil {
				return fmt.Errorf("failed to create bank account: %w", err)
			}
		}
	}

	return nil
}

// GetInstitutions retrieves available financial institutions for a country
func (s *gclService) GetInstitutions(ctx context.Context, countryCode string) ([]dto.Institution, error) {
	institutions, err := s.gclClient.GetInstitutions(countryCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get institutions: %w", err)
	}

	return institutions, nil
}

func (s *gclService) GetValidAccessToken() (string, error) {
	return s.gclClient.GetValidAccessToken()
}

func (s *gclService) RefreshTokenIfNeeded() error {
	_, err := s.gclClient.GetValidAccessToken()
	return err
}

func (s *gclService) GetTokenStatus() map[string]interface{} {
	return s.gclClient.GetTokenStatus()
}

func (s *gclService) ClearToken() {
	s.gclClient.ClearTokens()
}
