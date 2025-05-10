package server

import (
	"FinMa/types"
	"FinMa/utils"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var validate = validator.New()

// SignUpHandler is a handler that creates a new user.
// It expects a JSON object with the following fields:
// - email: the user's email address
// - password: the user's password
// - firstName: the user's first name
// - lastName: the user's last name
func (s *FiberServer) SignUpHandler(c *fiber.Ctx) error {
	var user types.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := validate.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Role = "user"

	if err := utils.ValidatePassword(user.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot hash password"})
	}
	user.Password = hashedPassword

	if err := s.db.CreateUser(user); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

func (s *FiberServer) LoginHandler(c *fiber.Ctx) error {
	var loginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.BodyParser(&loginRequest); err != nil {
		log.Error(fmt.Sprintf("error parsing login request: %s", err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := validate.Struct(loginRequest); err != nil {
		log.Error(fmt.Sprintf("error validating login request: %s", err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email or password format"})
	}

	user := s.db.GetUserByEmail(loginRequest.Email)

	if user.ID == uuid.Nil {
		log.Info(fmt.Sprintf("user not found: %s", loginRequest.Email))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}

	if err := utils.ComparePasswords(user.Password, loginRequest.Password); err != nil {
		log.Warn(fmt.Sprintf("invalid password for user: %s", loginRequest.Email))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid password"})
	}

	// Generate an access token
	payload := utils.Payload{
		UserID: user.ID,
		Email:  user.Email,
	}

	accessToken, err := utils.GenerateAccessToken(payload)

	if err != nil {
		log.Error(fmt.Sprintf("cannot generate access token: %s", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot generate access token"})
	}

	refreshToken, err := utils.GenerateRefreshToken(payload)

	if err != nil {
		log.Error(fmt.Sprintf("cannot generate refresh token: %s", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot generate refresh token"})
	}

	// return access token as a cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Minute * 5),
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
	})
	// Return user data without exposing sensitive information
	return c.JSON(fiber.Map{
		"id":    user.ID,
		"email": user.Email,
	})
}

func (s *FiberServer) RefreshHandler(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Refresh Token is missing",
		})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Verify the refresh token
	payload, err := utils.VerifyRefreshToken(req.RefreshToken)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid Refresh Token",
		})
	}

	existingUser := s.db.GetUserByEmail(payload.Email)

	if existingUser.ID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	accessToken, err := utils.GenerateAccessToken(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot generate access token",
		})
	}

	return c.JSON(fiber.Map{
		"access_token": accessToken,
	})
}

func (s *FiberServer) LogoutHandler(c *fiber.Ctx) error {
	// Clear the access token cookie
	c.ClearCookie("access_token")
	c.ClearCookie("refresh_token")
	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func (s *FiberServer) MeHandler(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*types.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user context",
		})
	}
	return c.JSON(fiber.Map{
		"message": "You are logged in",
		"user": fiber.Map{
			"id":        user.ID,
			"email":     user.Email,
			"role":      user.Role,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
	})
}
