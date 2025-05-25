package handlers

import (
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"FinMa/internal/domain"
	"FinMa/internal/service"
)

// GoCardlessHandler handles gocardless-related HTTP requests
type GoCardlessHandler struct {
	goCardlessService service.GoCardlessService
	validator         service.ValidatorService
}

// NewGoCardlessHandler creates a new gocardless handler
func NewGoCardlessHandler(goCardlessService service.GoCardlessService, validator service.ValidatorService) *GoCardlessHandler {
	return &GoCardlessHandler{
		goCardlessService: goCardlessService,
		validator:         validator,
	}
}

// Connect initiates the GoCardless OAuth2 consent flow
func (h *GoCardlessHandler) Connect(c *fiber.Ctx) error {
	// Get user from context (set by auth middleware)
	user, ok := c.Locals("user").(domain.User)
	if !ok {
		log.Error("Failed to get user from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	// Generate consent URL for GoCardless OAuth2 flow
	consentURL, err := h.goCardlessService.CreateConsentURL(c.Context(), user.ID)
	if err != nil {
		log.Error("Failed to create GoCardless consent URL", "error", err, "userID", user.ID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to initiate bank connection",
		})
	}

	return c.JSON(fiber.Map{
		"consent_url": consentURL,
		"message":     "Navigate to consent_url to authorize bank connection",
	})
}

// Callback handles the OAuth2 callback from GoCardless
func (h *GoCardlessHandler) Callback(c *fiber.Ctx) error {
	// Get authorization code from query params
	authCode := c.Query("code")
	if authCode == "" {
		log.Error("Missing authorization code in GoCardless callback")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Authorization code is required",
		})
	}

	// Check for error in callback
	if errorParam := c.Query("error"); errorParam != "" {
		errorDescription := c.Query("error_description")
		log.Error("GoCardless OAuth error", "error", errorParam, "description", errorDescription)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":             "Authorization failed",
			"error_description": errorDescription,
		})
	}

	// Get state parameter to verify request
	state := c.Query("state")
	if state == "" {
		log.Error("Missing state parameter in GoCardless callback")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid callback request",
		})
	}

	// Parse user ID from state (you should encode user ID in state during consent URL creation)
	userID, err := uuid.Parse(state)
	if err != nil {
		log.Error("Invalid state parameter", "state", state, "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid callback request",
		})
	}

	// Exchange authorization code for access token
	goCardlessItem, bankAccounts, err := h.goCardlessService.ExchangeAuthorizationCode(c.Context(), userID, authCode)
	if err != nil {
		log.Error("Failed to exchange authorization code", "error", err, "userID", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to complete bank connection",
		})
	}

	return c.JSON(fiber.Map{
		"message":         "Bank connection successful",
		"gocardless_item": goCardlessItem,
		"bank_accounts":   bankAccounts,
	})
}

// GetAccounts retrieves all bank accounts for the authenticated user
func (h *GoCardlessHandler) GetAccounts(c *fiber.Ctx) error {
	// Get user from context
	user, ok := c.Locals("user").(domain.User)
	if !ok {
		log.Error("Failed to get user from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	// Get user's bank accounts
	accounts, err := h.goCardlessService.GetUserBankAccounts(c.Context(), user.ID)
	if err != nil {
		log.Error("Failed to get user bank accounts", "error", err, "userID", user.ID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve bank accounts",
		})
	}

	return c.JSON(fiber.Map{
		"accounts": accounts,
	})
}

// RefreshAccounts manually refreshes account data from GoCardless
func (h *GoCardlessHandler) RefreshAccounts(c *fiber.Ctx) error {
	// Get user from context
	user, ok := c.Locals("user").(domain.User)
	if !ok {
		log.Error("Failed to get user from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	// Refresh account data
	err := h.goCardlessService.RefreshAccountData(c.Context(), user.ID)
	if err != nil {
		log.Error("Failed to refresh account data", "error", err, "userID", user.ID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to refresh account data",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Account data refreshed successfully",
	})
}

// DisconnectBank removes a bank connection
func (h *GoCardlessHandler) DisconnectBank(c *fiber.Ctx) error {
	// Get user from context
	user, ok := c.Locals("user").(domain.User)
	if !ok {
		log.Error("Failed to get user from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	// Get gocardless item ID from URL params
	goCardlessItemIDStr := c.Params("gocardless_item_id")
	if goCardlessItemIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "GoCardless item ID is required",
		})
	}

	goCardlessItemID, err := uuid.Parse(goCardlessItemIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid GoCardless item ID",
		})
	}

	// Disconnect the bank
	err = h.goCardlessService.DisconnectBank(c.Context(), user.ID, goCardlessItemID)
	if err != nil {
		log.Error("Failed to disconnect bank", "error", err, "userID", user.ID, "goCardlessItemID", goCardlessItemID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to disconnect bank",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Bank disconnected successfully",
	})
}

// GetTransactions retrieves transactions for a specific account
func (h *GoCardlessHandler) GetTransactions(c *fiber.Ctx) error {
	// Get user from context
	user, ok := c.Locals("user").(domain.User)
	if !ok {
		log.Error("Failed to get user from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	// Get account ID from URL params
	accountIDStr := c.Params("account_id")
	if accountIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Account ID is required",
		})
	}

	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account ID",
		})
	}

	// Get optional query parameters
	limit := c.QueryInt("limit", 100)
	offset := c.QueryInt("offset", 0)

	// Get transactions
	transactions, err := h.goCardlessService.GetAccountTransactions(c.Context(), user.ID, accountID, limit, offset)
	if err != nil {
		log.Error("Failed to get account transactions", "error", err, "userID", user.ID, "accountID", accountID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve transactions",
		})
	}

	return c.JSON(fiber.Map{
		"transactions": transactions,
		"limit":        limit,
		"offset":       offset,
	})
}
