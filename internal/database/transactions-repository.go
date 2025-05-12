package database

import (
	"FinMa/types"

	"github.com/charmbracelet/log"
)

func (s *service) CreateTransaction(transaction types.Transaction) error {
	s.db.Create(transaction)
	if s.db.Error != nil {
		return s.db.Error
	}
	return nil
}

func (s *service) GetTransactions(user types.User) []types.Transaction {
	var transactions []types.Transaction
	s.db.Where("user_id = ?", user.ID).Find(&transactions)

	if s.db.Error != nil {
		log.Error("Error fetching transactions: ", s.db.Error)
		return nil
	}
	return transactions
}

func (s *service) GetTransactionByID(id string) types.Transaction {
	var transaction types.Transaction
	s.db.First(&transaction, id)

	if s.db.Error != nil {
		log.Error("Error fetching transaction: ", s.db.Error)
		return types.Transaction{}
	}
	return transaction
}
