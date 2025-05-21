package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"FinMa/config"
	"FinMa/dto"
	"FinMa/internal/domain"
	"FinMa/internal/repository"
	"FinMa/utils"
)

// AuthService defines authentication service operations
type AuthService interface {
	Register(ctx context.Context, req dto.SignUpRequest) (domain.User, error)
	Login(ctx context.Context, req dto.LoginRequest) (domain.User, string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetUserByAccessToken(ctx context.Context, accessToken string) (domain.User, error)
	VerifyAccessToken(accessToken string) (Payload, error)
	VerifyRefreshToken(refreshToken string) (Payload, error)
	GenerateAccessToken(payload Payload) (string, error)
	GenerateRefreshToken(payload Payload) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
	config   *config.Config
}

// Payload represents the JWT payload data
type Payload struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, config *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		config:   config,
	}
}

// Register creates a new user
func (s *authService) Register(ctx context.Context, req dto.SignUpRequest) (domain.User, error) {
	// Check if user already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return domain.User{}, err
	}
	if exists {
		return domain.User{}, errors.New("user with this email already exists")
	}

	// Validate password strength
	if err := utils.ValidatePassword(req.Password); err != nil {
		return domain.User{}, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return domain.User{}, err
	}

	// Create user
	user := domain.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, &user); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

// Login authenticates a user and returns tokens
func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (domain.User, string, string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return domain.User{}, "", "", errors.New("invalid email or password")
	}

	// Verify password
	if err := utils.ComparePasswords(user.Password, req.Password); err != nil {
		return domain.User{}, "", "", errors.New("invalid email or password")
	}

	// Generate tokens
	payload := Payload{
		UserID: user.ID,
		Email:  user.Email,
	}

	accessToken, err := s.GenerateAccessToken(payload)
	if err != nil {
		return domain.User{}, "", "", err
	}

	refreshToken, err := s.GenerateRefreshToken(payload)
	if err != nil {
		return domain.User{}, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

// RefreshToken refreshes an access token
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// Verify refresh token
	payload, err := s.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	// Check if user exists
	_, err = s.userRepo.GetByEmail(ctx, payload.Email)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Generate new access token
	accessToken, err := s.GenerateAccessToken(payload)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// VerifyAccessToken verifies an access token and returns the user
func (s *authService) GetUserByAccessToken(ctx context.Context, accessToken string) (domain.User, error) {
	// Verify access token
	payload, err := s.VerifyAccessToken(accessToken)
	if err != nil {
		return domain.User{}, err
	}

	// Get user
	user, err := s.userRepo.GetByEmail(ctx, payload.Email)
	if err != nil {
		return domain.User{}, errors.New("user not found")
	}

	return user, nil
}

// GenerateAccessToken generates a new JWT access token.
// The payload is the data that will be stored in the token.
// The function returns the signed token as a string.
func (s *authService) GenerateAccessToken(payload Payload) (string, error) {
	if s.config.AccessTokenSecret == "" {
		return "", fmt.Errorf("access token secret is not set")
	}

	// Generate a new JWT token
	token := jwt.New()

	// Set the token claims
	token.Set("payload", payload)
	token.Set(jwt.IssuedAtKey, time.Now().Unix())
	token.Set(jwt.ExpirationKey, time.Now().Add(time.Minute*5).Unix())
	token.Set(jwt.IssuerKey, "FinMa")
	token.Set(jwt.SubjectKey, "access")
	token.Set(jwt.AudienceKey, "users")

	// Sign the token
	signedToken, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, []byte(s.config.AccessTokenSecret)))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return string(signedToken), nil
}

// GenerateRefreshToken generates a new JWT refresh token
func (s *authService) GenerateRefreshToken(payload Payload) (string, error) {
	if s.config.RefreshTokenSecret == "" {
		return "", fmt.Errorf("refresh token secret is not set")
	}

	// Generate a new JWT token
	token := jwt.New()

	// Set the token claims
	token.Set("payload", payload)
	token.Set(jwt.IssuedAtKey, time.Now().Unix())
	token.Set(jwt.ExpirationKey, time.Now().Add(time.Hour*24*7).Unix())
	token.Set(jwt.IssuerKey, "FinMa")
	token.Set(jwt.SubjectKey, "refresh")
	token.Set(jwt.AudienceKey, "users")

	// Sign the token
	signedToken, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, []byte(s.config.RefreshTokenSecret)))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return string(signedToken), nil
}

// VerifyAccessToken verifies the JWT access token
func (s *authService) VerifyAccessToken(accessToken string) (Payload, error) {
	if s.config.AccessTokenSecret == "" {
		return Payload{}, fmt.Errorf("access token secret is not set")
	}

	// Parse the token
	token, err := jwt.Parse(
		[]byte(accessToken),
		jwt.WithKey(jwa.HS256, []byte(s.config.AccessTokenSecret)),
		jwt.WithValidate(true),
	)
	if err != nil {
		return Payload{}, err
	}

	// Retrieve the payload from the token
	payload, ok := token.Get("payload")
	if !ok {
		return Payload{}, fmt.Errorf("failed to retrieve payload from token")
	}

	// Convert payload to Payload struct
	payloadMap := payload.(map[string]interface{})
	userID, err := uuid.Parse(payloadMap["user_id"].(string))
	if err != nil {
		return Payload{}, fmt.Errorf("failed to parse user_id: %w", err)
	}

	return Payload{
		UserID: userID,
		Email:  payloadMap["email"].(string),
	}, nil
}

// VerifyRefreshToken verifies the JWT refresh token
func (s *authService) VerifyRefreshToken(refreshToken string) (Payload, error) {
	if s.config.RefreshTokenSecret == "" {
		return Payload{}, fmt.Errorf("refresh token secret is not set")
	}

	// Parse the token
	token, err := jwt.Parse(
		[]byte(refreshToken),
		jwt.WithKey(jwa.HS256, []byte(s.config.RefreshTokenSecret)),
		jwt.WithValidate(true),
	)
	if err != nil {
		return Payload{}, err
	}

	// Retrieve the payload from the token
	payload, ok := token.Get("payload")
	if !ok {
		return Payload{}, fmt.Errorf("failed to retrieve payload from token")
	}

	// Convert payload to Payload struct
	payloadMap := payload.(map[string]interface{})
	userIDStr, ok := payloadMap["user_id"].(string)
	if !ok {
		return Payload{}, fmt.Errorf("user_id not found in token payload")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return Payload{}, fmt.Errorf("failed to parse user_id: %w", err)
	}

	email, ok := payloadMap["email"].(string)
	if !ok {
		return Payload{}, fmt.Errorf("email not found in token payload")
	}

	return Payload{
		UserID: userID,
		Email:  email,
	}, nil
}
