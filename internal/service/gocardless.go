package service

import (
	"context"

	"github.com/google/uuid"

	"FinMa/internal/domain"
	"FinMa/internal/repository"
)

type GoCardlessService interface {
	// OAuth2 flow
	CreateConsentURL(ctx context.Context, userID uuid.UUID) (string, error)
	ExchangeAuthorizationCode(ctx context.Context, userID uuid.UUID, authCode string) (domain.GoCardlessItem, []domain.BankAccount, error)

	// Account management
	GetUserBankAccounts(ctx context.Context, userID uuid.UUID) ([]domain.BankAccount, error)
	RefreshAccountData(ctx context.Context, userID uuid.UUID) error
	DisconnectBank(ctx context.Context, userID uuid.UUID, goCardlessItemID uuid.UUID) error

	// Transaction management
	GetAccountTransactions(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, limit, offset int) ([]domain.Transaction, error)
	SyncTransactions(ctx context.Context, userID uuid.UUID) error
}

type goCardlessService struct {
	goCardlessItemRepo repository.GoCardlessItemRepository
	bankAccountRepo    repository.BankAccountRepository
	userRepo           repository.UserRepository
}

// NewGoCardlessService creates a new GoCardless service
func NewGoCardlessService(
	goCardlessItemRepo repository.GoCardlessItemRepository,
	bankAccountRepo repository.BankAccountRepository,
	userRepo repository.UserRepository,
) GoCardlessService {
	return &goCardlessService{
		goCardlessItemRepo: goCardlessItemRepo,
		bankAccountRepo:    bankAccountRepo,
		userRepo:           userRepo,
	}
}

// GoCardlessClient interface for GoCardless API operations
type GoCardlessClient interface {
	CreateConsentURL(userID uuid.UUID, redirectURI string) (string, error)
	ExchangeAuthorizationCode(authCode string) (GoCardlessTokenResponse, error)
	GetAccounts(accessToken string) ([]GoCardlessAccount, error)
	GetTransactions(accessToken, accountID string, limit, offset int) ([]GoCardlessTransaction, error)
	RefreshAccessToken(refreshToken string) (GoCardlessTokenResponse, error)
}

// GoCardless API response types
type GoCardlessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type GoCardlessAccount struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Number           string  `json:"number"`
	Type             string  `json:"type"`
	Currency         string  `json:"currency"`
	InstitutionName  string  `json:"institution_name"`
	BalanceAvailable float64 `json:"balance_available"`
	BalanceCurrent   float64 `json:"balance_current"`
}

type GoCardlessTransaction struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	Category    string  `json:"category"`
	Type        string  `json:"type"`
}

// Implementation methods would go here...
func (s *goCardlessService) CreateConsentURL(ctx context.Context, userID uuid.UUID) (string, error) {
	return "", nil
}

func (s *goCardlessService) ExchangeAuthorizationCode(ctx context.Context, userID uuid.UUID, authCode string) (domain.GoCardlessItem, []domain.BankAccount, error) {
	// Implementation here
	// 1. Exchange auth code for tokens
	// 2. Save GoCardlessItem to database
	// 3. Fetch and save bank accounts
	// 4. Return created items
	return domain.GoCardlessItem{}, []domain.BankAccount{}, nil
}

func (s *goCardlessService) GetUserBankAccounts(ctx context.Context, userID uuid.UUID) ([]domain.BankAccount, error) {
	return s.bankAccountRepo.GetByUserID(ctx, userID)
}

func (s *goCardlessService) RefreshAccountData(ctx context.Context, userID uuid.UUID) error {
	return nil
}

func (s *goCardlessService) DisconnectBank(ctx context.Context, userID uuid.UUID, goCardlessItemID uuid.UUID) error {
	return s.goCardlessItemRepo.Delete(ctx, goCardlessItemID)
}

func (s *goCardlessService) GetAccountTransactions(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, limit, offset int) ([]domain.Transaction, error) {
	return []domain.Transaction{}, nil
}

func (s *goCardlessService) SyncTransactions(ctx context.Context, userID uuid.UUID) error {
	return nil
}
