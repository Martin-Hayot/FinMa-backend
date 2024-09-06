package server

import (
	"FinMa/internal/database"
	"FinMa/utils"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var validate = validator.New()

func (s *FiberServer) RegisterFiberRoutes() {
	api := s.Group("/api")
	auth := api.Group("/auth")
	api.Get("/", s.HelloWorldHandler)
	api.Get("/health", s.healthHandler)
	auth.Post("/signup", s.SignUpHandler)
	auth.Post("/login", s.SignUpHandler)
	auth.Post("/logout", s.SignUpHandler)
	api.Post("/users", s.SignUpHandler)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}

func (s *FiberServer) GetUsersHandler(c *fiber.Ctx) error {
	db := s.db.GetUsers()

	return c.JSON(db)
}

// SignUpHandler is a handler that creates a new user.
// It expects a JSON object with the following fields:
// - email: the user's email address
// - password: the user's password
// - first_name: the user's first name
// - last_name: the user's last name
func (s *FiberServer) SignUpHandler(c *fiber.Ctx) error {
	var user database.User

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

	if err := utils.ValidatePassword(user.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot hash password"})
	}
	user.Password = hashedPassword

	if err := s.db.CreateUser(user); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}
