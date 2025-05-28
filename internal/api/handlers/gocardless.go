package handlers

import (
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"

	"FinMa/internal/service"
)

// GclHandler handles gocardless-related HTTP requests
type GclHandler struct {
	goCardlessService service.GclService
	validator         service.ValidatorService
}

// NewGclHandler creates a new gocardless handler
func NewGclHandler(goCardlessService service.GclService, validator service.ValidatorService) *GclHandler {
	return &GclHandler{
		goCardlessService: goCardlessService,
		validator:         validator,
	}
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
