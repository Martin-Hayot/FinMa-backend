package server

import (
	"FinMa/utils"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// Authorize is a middleware that checks if the user is authenticated and has the correct role.
// If the user is not authenticated or doesn't have the correct role, it returns a 401 Unauthorized error.
func (s *FiberServer) Authorize(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Warn("Authorization header is missing")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		auth := strings.Fields(authHeader)
		if len(auth) != 2 || auth[0] != "Bearer" {
			log.Warn("Invalid Authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Check if the token is valid
		token := auth[1]
		payload, err := utils.VerifyAccessToken(token)
		if err != nil {
			if err.Error() == jwt.ErrTokenExpired().Error() {
				log.Warn("Access token expired:", err)
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Token expired, please refresh",
				})
			}
			log.Warn("Invalid access token:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		existingUser := s.db.GetUserByEmail(payload.Email)

		if existingUser.ID == uuid.Nil {
			log.Warn("User not found")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Check if the user has the correct role
		if !utils.HasRole(existingUser.Role, allowedRoles) {
			log.Warnf("User %s does not have the required role", existingUser.Email)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden: You do not have permission to access this resource",
			})
		}

		// Store the user in the context
		c.Locals("user", existingUser)

		// Continue to the next middleware
		return c.Next()
	}
}
