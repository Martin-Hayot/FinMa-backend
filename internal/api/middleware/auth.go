package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"FinMa/internal/service"
)

// AuthMiddleware creates middleware for authentication validation
func AuthMiddleware(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the access token from cookies or Authorization header
		var accessToken string

		// Try cookie first
		accessToken = c.Cookies("access_token")

		// If not in cookie, try Authorization header
		if accessToken == "" {
			authHeader := c.Get("Authorization")
			if authHeader == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Missing authentication token",
				})
			}

			// Check for Bearer token format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid authorization format",
				})
			}

			accessToken = parts[1]
		}

		// Verify the token
		user, err := authService.GetUserByAccessToken(c.Context(), accessToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Add user to context
		c.Locals("user", user)

		return c.Next()
	}
}
