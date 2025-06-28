package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"FinMa/dto"
	"FinMa/internal/domain"
	"FinMa/internal/repository"
	"FinMa/pkg/gocardless"

	"github.com/google/uuid"
)

type GclService interface {
	// Token management
	GetValidAccessToken(ctx context.Context) (string, error)
	RefreshTokenIfNeeded(ctx context.Context) error
	GetTokenStatus() map[string]interface{}
	ClearToken()

	// LinkAccount initiates the linking of a bank account for a user with a specific institution
	LinkAccount(ctx context.Context, userID uuid.UUID, institutionID, redirectURL string) (*dto.LinkAccountResponse, error)

	// SyncRequisition syncs an existing requisition for a user
	SyncRequisition(ctx context.Context, requisitionReference string, userID uuid.UUID) (*dto.GoCardlessUpdateRequisitionResponse, error)

	// Institutions
	GetInstitutions(ctx context.Context, countryCode string) ([]dto.Institution, error)

	// Account Data
	GetAccountDetails(ctx context.Context, accountID string) (*dto.AccountDetails, error)
	GetAccountBalances(ctx context.Context, accountID string) (*dto.AccountBalances, error)
	GetAccountTransactions(ctx context.Context, accountID string) (*dto.AccountTransactions, error)
}

type gclService struct {
	bankAccountRepo repository.BankAccountRepository
	userRepo        repository.UserRepository
	requisitionRepo repository.RequisitionRepository
	transactionRepo repository.TransactionRepository
	gclClient       *gocardless.Client
}

// NewGclService creates a new GoCardless service
func NewGclService(
	bankAccountRepo repository.BankAccountRepository,
	userRepo repository.UserRepository,
	requisitionRepo repository.RequisitionRepository,
	transactionRepo repository.TransactionRepository,
	gclClient *gocardless.Client,
) GclService {
	return &gclService{
		bankAccountRepo: bankAccountRepo,
		requisitionRepo: requisitionRepo,
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
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
	response, err := s.gclClient.CreateRequisition(ctx, userID, institutionID, redirectURL)
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

	if requisition.Status == "LN" {
		// If the requisition is already linked, return early
		return &dto.GoCardlessUpdateRequisitionResponse{
			Status:        requisition.Status,
			InstitutionID: requisition.InstitutionID,
			Reference:     requisition.Reference,
		}, nil
	}

	// Call the GoCardless client to update the requisition
	response, err := s.gclClient.GetRequisition(ctx, userID, requisition.ID)
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
		// Fetch account details and balances
		accountDetails, err := s.gclClient.GetAccountDetails(ctx, accountID)
		if err != nil {
			return fmt.Errorf("failed to get account details for %s: %w", accountID, err)
		}

		balances, err := s.gclClient.GetAccountBalances(ctx, accountID)
		if err != nil {
			return fmt.Errorf("failed to get account balances for %s: %w", accountID, err)
		}

		var balanceAvailable, balanceCurrent float64
		for _, balance := range balances.Balances {
			if balance.BalanceType == "interimAvailable" {
				balanceAvailable, _ = strconv.ParseFloat(balance.BalanceAmount.Amount, 64)
			}
			if balance.BalanceType == "interimBooked" {
				balanceCurrent, _ = strconv.ParseFloat(balance.BalanceAmount.Amount, 64)
			}
		}

		// Check if the account already exists in our database
		existingAccount, err := s.bankAccountRepo.GetByAccountID(ctx, accountID)
		if err != nil {
			// If the error is a 'not found' error, it means we can create the account.
			if repository.IsNotFoundError(err) {
				// Create new bank account record
				bankAccount := &domain.BankAccount{
					ID:               uuid.New(),
					AccountID:        accountID,
					Name:             accountDetails.Account.Name,
					Type:             accountDetails.Account.Product,
					Currency:         accountDetails.Account.Currency,
					InstitutionName:  accountDetails.Account.InstitutionName,
					IBAN:             accountDetails.Account.IBAN,
					UserID:           userID,
					RequisitionID:    requisitionID,
					BalanceAvailable: balanceAvailable,
					BalanceCurrent:   balanceCurrent,
				}

				err = s.bankAccountRepo.Create(ctx, bankAccount)
				if err != nil {
					return fmt.Errorf("failed to create bank account: %w", err)
				}
				existingAccount = bankAccount // Set existingAccount for transaction processing
			} else {
				// Any other error is a real problem.
				return fmt.Errorf("failed to check existing account: %w", err)
			}
		} else {
			// If account exists, update its balances
			existingAccount.BalanceAvailable = balanceAvailable
			existingAccount.BalanceCurrent = balanceCurrent
			err = s.bankAccountRepo.Update(ctx, existingAccount)
			if err != nil {
				return fmt.Errorf("failed to update existing bank account balances: %w", err)
			}
		}

		// Process transactions for the account
		err = s.processTransactionsForAccount(ctx, accountID, existingAccount.ID, userID)
		if err != nil {
			return fmt.Errorf("failed to process transactions for account %s: %w", accountID, err)
		}
	}

	return nil
}

func (s *gclService) processTransactionsForAccount(ctx context.Context, accountID string, bankAccountID uuid.UUID, userID uuid.UUID) error {
	transactions, err := s.gclClient.GetAccountTransactions(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to get transactions from GoCardless for account %s: %w", accountID, err)
	}

	var newTransactions []*domain.Transaction
	for _, tx := range transactions.Transactions.Booked {
		// Check if transaction already exists to avoid duplicates
		exists, err := s.transactionRepo.ExistsByTransactionID(ctx, tx.TransactionID)
		if err != nil {
			return fmt.Errorf("failed to check existence of transaction %s: %w", tx.TransactionID, err)
		}
		if exists {
			continue // Skip existing transactions
		}

		amount, _ := strconv.ParseFloat(tx.TransactionAmount.Amount, 64)
		bookingDate, _ := time.Parse("2006-01-02", tx.BookingDate)

		newTransactions = append(newTransactions, &domain.Transaction{
			ID:            uuid.New(),
			Description:   tx.RemittanceInformation,
			Amount:        amount,
			Date:          bookingDate,
			Type:          "",    // You might need to infer this from category or other logic
			IsRecurring:   false, // You might need to infer this
			UserID:        userID,
			BankAccountID: bankAccountID,
		})
	}

	if len(newTransactions) > 0 {
		err = s.transactionRepo.CreateInBatches(ctx, newTransactions)
		if err != nil {
			return fmt.Errorf("failed to save new transactions: %w", err)
		}
	}

	return nil
}

// GetInstitutions retrieves available financial institutions for a country
func (s *gclService) GetInstitutions(ctx context.Context, countryCode string) ([]dto.Institution, error) {
	institutions, err := s.gclClient.GetInstitutions(ctx, countryCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get institutions: %w", err)
	}

	return institutions, nil
}

func (s *gclService) GetAccountDetails(ctx context.Context, accountID string) (*dto.AccountDetails, error) {
	return s.gclClient.GetAccountDetails(ctx, accountID)
}

func (s *gclService) GetAccountBalances(ctx context.Context, accountID string) (*dto.AccountBalances, error) {
	return s.gclClient.GetAccountBalances(ctx, accountID)
}

func (s *gclService) GetAccountTransactions(ctx context.Context, accountID string) (*dto.AccountTransactions, error) {
	return s.gclClient.GetAccountTransactions(ctx, accountID)
}

func (s *gclService) GetValidAccessToken(ctx context.Context) (string, error) {
	return s.gclClient.GetValidAccessToken(ctx)
}

func (s *gclService) RefreshTokenIfNeeded(ctx context.Context) error {
	_, err := s.gclClient.GetValidAccessToken(ctx)
	return err
}

func (s *gclService) GetTokenStatus() map[string]interface{} {
	return s.gclClient.GetTokenStatus()
}

func (s *gclService) ClearToken() {
	s.gclClient.ClearTokens()
}
