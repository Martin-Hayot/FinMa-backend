package database

import (
	"FinMa/types"
	"fmt"

	"github.com/google/uuid"
)

func (s *service) GetUsers() []types.User {
	var users []types.User
	s.db.Find(&users)
	return users
}

func (s *service) GetUser(id int) types.User {
	var user types.User
	s.db.First(&user, id)
	return user
}

func (s *service) CreateUser(user types.User) error {
	// Check if user already exists
	var existingUser types.User
	result := s.db.Where("email = ?", user.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}
	// Create the new user
	s.db.Create(&user)
	return nil
}

func (s *service) GetUserByEmail(email string) types.User {
	var user types.User
	s.db.Where("email = ?", email).First(&user)
	return user
}

func (s *service) GetUserByID(id uuid.UUID) types.User {
	var user types.User
	s.db.Where("id = ?", id).First(&user)
	return user
}
