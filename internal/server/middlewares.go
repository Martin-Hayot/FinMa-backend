package server

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Authorize is a middleware that checks if the user is authenticated.
// If the user is not authenticated, it returns a 401 Unauthorized error.
// The header must contain an Authorization token.
// Example: Authorization: Bearer <token>
func Authorize(c *fiber.Ctx) error {
	// Check if the user is authenticated
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Check if the token is valid
	auth := strings.Fields(authHeader)
	if len(auth) != 2 || auth[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Continue to the next middleware
	return c.Next()
}
