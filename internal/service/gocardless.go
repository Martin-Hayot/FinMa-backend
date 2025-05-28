package service

import (
	"context"
	"fmt"

	"FinMa/dto"
	"FinMa/internal/repository"
	"FinMa/pkg/gocardless"
)

type GclService interface {
	// Token management
	GetValidAccessToken() (string, error)
	RefreshTokenIfNeeded() error
	GetTokenStatus() map[string]interface{}
	ClearToken()
	// Institutions
	GetInstitutions(ctx context.Context, countryCode string) ([]dto.Institution, error)
}

type gclService struct {
	gclItemRepo     repository.GclItemRepository
	bankAccountRepo repository.BankAccountRepository
	userRepo        repository.UserRepository
	gclClient       *gocardless.Client
}

// NewGclService creates a new GoCardless service
func NewGclService(
	gclItemRepo repository.GclItemRepository,
	bankAccountRepo repository.BankAccountRepository,
	userRepo repository.UserRepository,
	gclClient *gocardless.Client,
) GclService {
	return &gclService{
		gclItemRepo:     gclItemRepo,
		bankAccountRepo: bankAccountRepo,
		userRepo:        userRepo,
		gclClient:       gclClient,
	}
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
