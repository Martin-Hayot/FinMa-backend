package service

import (
	"context"
	"fmt"

	"FinMa/dto"
	"FinMa/internal/repository"

	"github.com/google/uuid"
)

type BankAccountService interface {
	GetBankAccountsForUser(ctx context.Context, userID uuid.UUID) ([]dto.BankAccountResponse, error)
}

type bankAccountService struct {
	bankAccountRepo repository.BankAccountRepository
}

func NewBankAccountService(bankAccountRepo repository.BankAccountRepository) BankAccountService {
	return &bankAccountService{
		bankAccountRepo: bankAccountRepo,
	}
}

func (s *bankAccountService) GetBankAccountsForUser(ctx context.Context, userID uuid.UUID) ([]dto.BankAccountResponse, error) {
	accounts, err := s.bankAccountRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bank accounts for user %s: %w", userID, err)
	}

	var response []dto.BankAccountResponse
	for _, acc := range accounts {
		response = append(response, dto.BankAccountResponse{
			ID:               acc.ID,
			AccountID:        acc.AccountID,
			Name:             acc.Name,
			Type:             acc.Type,
			Currency:         acc.Currency,
			InstitutionName:  acc.InstitutionName,
			BalanceAvailable: acc.BalanceAvailable,
			BalanceCurrent:   acc.BalanceCurrent,
			IBAN:             acc.IBAN,
			CreatedAt:        acc.CreatedAt,
			UpdatedAt:        acc.UpdatedAt,
		})
	}

	return response, nil
}
