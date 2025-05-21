package handlers

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"

	"FinMa/dto"
	"FinMa/internal/domain"
	"FinMa/internal/service"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService service.AuthService
	validator   service.ValidatorService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService, validator service.ValidatorService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
	}
}

// SignUp handles user registration
func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
	var req dto.SignUpRequest
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

	// Register user
	user, err := h.authService.Register(ctx, req)
	if err != nil {
		log.Error("Failed to register user", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}

	// Return user response
	return c.Status(fiber.StatusCreated).JSON(dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request body
	var req dto.LoginRequest
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

	// Authenticate user
	user, accessToken, refreshToken, err := h.authService.Login(ctx, req)
	if err != nil {
		log.Error("Failed to login", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Set cookies
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Minute * 15),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	// Return user info
	return c.JSON(dto.LoginResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}

// Refresh handles token refreshing
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get refresh token from cookies
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		// Try to get from request body
		var req dto.RefreshTokenRequest
		if err := c.BodyParser(&req); err != nil || req.RefreshToken == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Refresh token is required",
			})
		}
		refreshToken = req.RefreshToken
	}

	// Refresh token
	accessToken, err := h.authService.RefreshToken(ctx, refreshToken)
	if err != nil {
		log.Error("Failed to refresh token", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid refresh token",
		})
	}

	// Set new access token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Minute * 15),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return c.JSON(dto.RefreshTokenResponse{
		AccessToken: accessToken,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Clear cookies
	c.ClearCookie("access_token")
	c.ClearCookie("refresh_token")

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

// Me returns the current user's info
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(domain.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	return c.JSON(dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	})
}
