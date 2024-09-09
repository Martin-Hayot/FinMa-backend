package handlers

import (
	"FinMa/constants"
	"FinMa/internal/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateTransaction(c *fiber.Ctx) error {
	db := database.Get()

	type CreateTransactionRequest struct {
		Category      string    `json:"category"`
		Amount        float64   `json:"amount"`
		Date          time.Time `json:"date"`
		Type          string    `json:"type"` // income/expense
		IsRecurring   bool      `json:"is_recurring"`
		Description   string    `json:"description"`
		BankAccountID uuid.UUID `json:"bank_account_id"`
	}

	var body CreateTransactionRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// validate the types
	validType := false
	for _, t := range constants.GetTransactionTypes() {
		if t == body.Type {
			validType = true
			break
		}
	}
	if !validType {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction type",
		})
	}

	// validate the category
	validCategory := false
	for _, cat := range constants.GetTransactionCategories() {
		if cat == body.Category {
			validCategory = true
			break
		}
	}
	if !validCategory {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction category",
		})
	}

	transaction := &database.Transaction{
		Category:      body.Category,
		Amount:        body.Amount,
		Date:          body.Date,
		Type:          body.Type,
		IsRecurring:   body.IsRecurring,
		Description:   body.Description,
		BankAccountID: body.BankAccountID,
	}

	if err := db.CreateTransaction(transaction); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create transaction",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(transaction)
}

func GetTransactions(c *fiber.Ctx) error {
	db := database.Get()
	user := c.Locals("user").(*database.User)
	transactions := db.GetTransactions(user)

	return c.JSON(transactions)
}

func GetTransactionByID(c *fiber.Ctx) error {
	db := database.Get()
	id := c.Params("id")
	transaction := db.GetTransactionByID(id)

	if transaction.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	return c.JSON(transaction)
}
