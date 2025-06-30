package handlers

import (
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"FinMa/dto"
	"FinMa/internal/domain"
	"FinMa/internal/service"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService service.UserService
	validator   service.ValidatorService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService, validator service.ValidatorService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator,
	}
}

// GetProfile handles retrieving the authenticated user's profile
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		log.Error("Failed to get user ID from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	// Get user profile
	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		log.Error("Failed to get user profile", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user profile",
		})
	}

	return c.JSON(user)
}

// DeleteAccount handles deleting the authenticated user's account
func (h *UserHandler) DeleteAccount(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from context (set by auth middleware)
	user, ok := c.Locals("user").(domain.User)
	if !ok {
		log.Error("Failed to get user from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}
	userID := user.ID

	// Parse request body
	var req struct {
		Password string `json:"password" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Delete account
	err := h.userService.DeleteAccount(ctx, userID, req.Password)
	if err != nil {
		if err.Error() == "incorrect password" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Password is incorrect",
			})
		}
		log.Error("Failed to delete account", "error", err, "userID", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete account",
		})
	}

	// Clear authentication cookies
	c.ClearCookie("access_token")
	c.ClearCookie("refresh_token")

	return c.JSON(fiber.Map{
		"message": "Account deleted successfully",
	})
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user from context (set by auth middleware)
	user, ok := c.Locals("user").(domain.User)
	if !ok {
		log.Error("Failed to get user from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	// Parse request body
	var req dto.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Update user
	updatedUser, err := h.userService.UpdateProfile(ctx, user, req)
	if err != nil {
		log.Error("Failed to update user", "error", err, "userID", user.ID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.JSON(updatedUser)
}
