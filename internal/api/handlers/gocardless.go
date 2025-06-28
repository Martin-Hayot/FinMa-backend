package handlers

import (
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"

	"FinMa/config"
	"FinMa/dto"
	"FinMa/internal/domain"
	"FinMa/internal/service"
)

// GclHandler handles gocardless-related HTTP requests
type GclHandler struct {
	goCardlessService service.GclService
	cfg               *config.Config
	validator         service.ValidatorService
}

// NewGclHandler creates a new gocardless handler
func NewGclHandler(goCardlessService service.GclService, validator service.ValidatorService, cfg *config.Config) *GclHandler {
	return &GclHandler{
		goCardlessService: goCardlessService,
		cfg:               cfg,
		validator:         validator,
	}
}

func (h *GclHandler) LinkAccount(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	user, ok := c.Locals("user").(domain.User)
	if !ok {
		log.Error("Failed to get user ID from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	var req dto.LinkAccountRequest

	// Validate request body
	if err := c.BodyParser(&req); err != nil {
		log.Error("Failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Call GoCardless service to create requisition
	requisition, err := h.goCardlessService.LinkAccount(c.Context(), user.ID, req.InstitutionID, h.cfg.GoCardless.RedirectURL)
	if err != nil {
		log.Error("Failed to create requisition", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create requisition",
		})
	}

	return c.JSON(requisition)
}

func (h *GclHandler) SyncRequisition(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	user, ok := c.Locals("user").(domain.User)
	if !ok {
		log.Error("Failed to get user ID from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}
	// Get requisition ID from URL params
	requisitionReference := c.Params("id")
	if requisitionReference == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Requisition ID is required",
		})
	}

	// Call GoCardless service to update requisition
	response, err := h.goCardlessService.SyncRequisition(c.Context(), requisitionReference, user.ID)
	if err != nil {
		log.Error("Failed to sync requisition", "error", err, "requisitionReference", requisitionReference)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to sync requisition",
		})
	}

	return c.JSON(response)
}

// GetInstitutions retrieves available financial institutions for a country
func (h *GclHandler) GetInstitutions(c *fiber.Ctx) error {
	// Get country code from URL params
	countryCode := c.Params("country_code")
	if countryCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Country code is required",
		})
	}

	// Validate country code format (should be 2 characters)
	if len(countryCode) != 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Country code must be 2 characters (e.g., 'gb', 'fr')",
		})
	}

	// Get institutions from GoCardless service
	institutions, err := h.goCardlessService.GetInstitutions(c.Context(), countryCode)
	if err != nil {
		log.Error("Failed to get institutions", "error", err, "countryCode", countryCode)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve institutions",
		})
	}

	return c.JSON(fiber.Map{
		"institutions": institutions,
		"country_code": countryCode,
	})
}

func (h *GclHandler) GetTokenStatus(c *fiber.Ctx) error {
	status := h.goCardlessService.GetTokenStatus()
	return c.JSON(status)
}
