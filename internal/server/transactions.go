package server

import (
	"FinMa/constants"
	"FinMa/types"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (s *FiberServer) CreateTransaction(c *fiber.Ctx) error {
	type CreateTransactionRequest struct {
		Category      string    `json:"category"`
		Amount        float64   `json:"amount"`
		Date          string    `json:"date"` // Change to string for custom parsing
		Type          string    `json:"type"` // income/expense
		IsRecurring   bool      `json:"is_recurring"`
		Description   string    `json:"description"`
		BankAccountID uuid.UUID `json:"bank_account_id"`
	}

	var body CreateTransactionRequest
	if err := c.BodyParser(&body); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate the date format
	parsedDate, err := time.Parse(time.RFC3339, body.Date)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format",
		})
	}

	// Validate the types
	validType := false
	for _, t := range constants.GetTransactionTypes() {
		if t == body.Type {
			validType = true
			break
		}
	}
	if !validType {
		log.Error("Invalid transaction type")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction type",
		})
	}

	// Validate the category
	validCategory := false
	for _, cat := range constants.GetTransactionCategories() {
		if cat == body.Category {
			validCategory = true
			break
		}
	}
	if !validCategory {
		log.Error("Invalid transaction category")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction category",
		})
	}

	user := c.Locals("user").(types.User)

	if user.ID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	transaction := types.Transaction{
		Category:      body.Category,
		Amount:        body.Amount,
		Date:          parsedDate,
		Type:          body.Type,
		IsRecurring:   body.IsRecurring,
		Description:   body.Description,
		BankAccountID: body.BankAccountID,
		User:          user,
	}

	if err := s.db.CreateTransaction(transaction); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create transaction",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(transaction)
}

func (s *FiberServer) GetTransactions(c *fiber.Ctx) error {
	user := c.Locals("user").(types.User)
	transactions := s.db.GetTransactions(user)

	return c.JSON(transactions)
}

func (s *FiberServer) GetTransactionByID(c *fiber.Ctx) error {
	id := c.Params("id")
	transaction := s.db.GetTransactionByID(id)

	if transaction.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	return c.JSON(transaction)
}
