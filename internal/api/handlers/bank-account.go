package handlers

import (
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"

	"FinMa/internal/domain"
	"FinMa/internal/service"
)

// BankAccountHandler handles bank account related HTTP requests
type BankAccountHandler struct {
	bankAccountService service.BankAccountService
}

// NewBankAccountHandler creates a new bank account handler
func NewBankAccountHandler(bankAccountService service.BankAccountService) *BankAccountHandler {
	return &BankAccountHandler{
		bankAccountService: bankAccountService,
	}
}

// GetAccounts retrieves all bank accounts for the authenticated user
func (h *BankAccountHandler) GetAccounts(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	user, ok := c.Locals("user").(domain.User)
	if !ok {
		log.Error("Failed to get user ID from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	accounts, err := h.bankAccountService.GetBankAccountsForUser(c.Context(), user.ID)
	if err != nil {
		log.Error("Failed to get bank accounts", "error", err, "userID", user.ID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve bank accounts",
		})
	}

	return c.JSON(accounts)
}