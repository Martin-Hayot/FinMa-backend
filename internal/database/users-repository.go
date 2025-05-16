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

// CreateEmailVerificationToken creates a new email verification token in the database.
func (s *service) CreateEmailVerificationToken(token types.EmailVerificationToken) error {
    return s.db.Create(&token).Error
}

// GetEmailVerificationToken retrieves an email verification token by its value.
func (s *service) GetEmailVerificationToken(token string) types.EmailVerificationToken {
    var verificationToken types.EmailVerificationToken
    s.db.Where("token = ?", token).First(&verificationToken)
    return verificationToken
}

// DeleteEmailVerificationToken deletes an email verification token by its ID.
func (s *service) DeleteEmailVerificationToken(id uuid.UUID) error {
    return s.db.Delete(&types.EmailVerificationToken{}, id).Error
}

// UpdateUser updates the user record in the database.
func (s *service) UpdateUser(user types.User) error {
    return s.db.Save(&user).Error
}
// DeleteEmailVerificationTokenByUserID supprime tous les tokens de vérification associés à un utilisateur.
func (s *service) DeleteEmailVerificationTokenByUserID(userID uuid.UUID) error {
    return s.db.Where("user_id = ?", userID).Delete(&types.EmailVerificationToken{}).Error
}