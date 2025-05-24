package handlers

import (
	"FinMa/dto"
	"FinMa/internal/domain"
	"FinMa/internal/service"
	"FinMa/plaid"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type PlaidHandler struct {
	PlaidClient        plaid.Client
	Validator          service.ValidatorService
	UserService        service.UserService
	BankAccountService service.BankAccountService
	PlaidItemService   service.PlaidItemService
}

// NewPlaidHandler creates a new Plaid handler
func NewPlaidHandler(
	plaidClient plaid.Client,
	validator service.ValidatorService,
	userService service.UserService,
	bankAccountService service.BankAccountService,
	plaidItemService service.PlaidItemService,
) *PlaidHandler {
	return &PlaidHandler{
		PlaidClient:        plaidClient,
		Validator:          validator,
		UserService:        userService,
		BankAccountService: bankAccountService,
		PlaidItemService:   plaidItemService,
	}
}

// CreateLinkToken handles the creation of a Plaid link token
func (h *PlaidHandler) CreateLinkToken(c *fiber.Ctx) error {
	ctx := c.Context()
	// Get user from context (set by auth middleware)
	user, ok := c.Locals("user").(domain.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	log.Debug("Creating Plaid link token for user ID:", user.ID)
	// Create link token
	linkToken, err := h.PlaidClient.CreateLinkToken(ctx, user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create link token",
		})
	}
	// Return the link token to the client
	return c.JSON(fiber.Map{
		"link_token": linkToken,
	})
}

// ExchangePublicToken handles the exchange of a public token for an access token
func (h *PlaidHandler) ExchangePublicToken(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user from context
	user, ok := c.Locals("user").(domain.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	// Get public token from request body
	var requestBody struct {
		PublicToken string `json:"public_token"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	publicToken := requestBody.PublicToken
	if publicToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Public token is required",
		})
	}

	// 1. Exchange the public token for an access token
	accessToken, itemID, err := h.PlaidClient.ExchangePublicToken(ctx, publicToken)
	if err != nil {
		log.Error("Failed to exchange public token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to exchange public token",
		})
	}

	// 2. Save PlaidItem (and link to User)
	plaidItem, err := h.PlaidItemService.CreatePlaidItem(ctx, user.ID, accessToken, itemID)
	if err != nil {
		log.Error("Failed to create Plaid item:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save Plaid item",
		})
	}

	// 3. Fetch accounts from Plaid
	h.PlaidClient.AccessToken = accessToken
	accountsResp, err := h.PlaidClient.GetAccounts(ctx)
	if err != nil {
		log.Error("Failed to get accounts from Plaid:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve accounts",
		})
	}

	// 4. Save accounts under BankAccount, linked to User and PlaidItem
	var responseAccounts []dto.BankAccountResponse
	for _, account := range accountsResp.GetAccounts() {
		bankAccount := domain.BankAccount{
			AccountID:        account.GetAccountId(),
			AccountName:      account.GetName(),
			OfficialName:     account.GetOfficialName(),
			AccountType:      string(account.GetType()),
			CurrentBalance:   *account.GetBalances().Current.Get(),
			AvailableBalance: *account.GetBalances().Available.Get(),
			Mask:             account.GetMask(),
		}

		createdAccount, err := h.BankAccountService.CreateBankAccount(ctx, user.ID, plaidItem.ID, bankAccount)
		if err != nil {
			log.Error("Failed to create bank account:", err)
			continue // Continue with other accounts
		}
		// Convert domain.BankAccount to dto.BankAccountResponse before appending
		responseAccount := dto.BankAccountResponse{
			AccountID:    createdAccount.AccountID,
			AccountName:  createdAccount.AccountName,
			OfficialName: createdAccount.OfficialName,
			AccountType:  createdAccount.AccountType,
			Mask:         createdAccount.Mask,
		}
		responseAccounts = append(responseAccounts, responseAccount)
	}

	// 5. Return success with account summary
	return c.JSON(fiber.Map{
		"success":        true,
		"item_id":        itemID,
		"accounts_count": len(responseAccounts),
		"accounts":       responseAccounts,
		"message":        "Successfully connected bank accounts",
	})
}
