package dto

import "github.com/google/uuid"

// SignUpRequest represents the data needed for user registration
type SignUpRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

// LoginRequest represents the data needed for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the data returned after successful login
type LoginResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

// RefreshTokenRequest represents the request to refresh an access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse represents the response with a new access token
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}
