package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"FinMa/dto"
	"FinMa/internal/api/middleware"
	"FinMa/internal/domain"
	"FinMa/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockAuthService is a mock implementation of the AuthService for testing.
type MockAuthService struct {
	RegisterFunc             func(ctx context.Context, req dto.SignUpRequest) (domain.User, error)
	LoginFunc                func(ctx context.Context, req dto.LoginRequest) (domain.User, string, string, error)
	RefreshTokenFunc         func(ctx context.Context, refreshToken string) (string, error)
	GetUserByAccessTokenFunc func(ctx context.Context, accessToken string) (domain.User, error)
	VerifyAccessTokenFunc    func(accessToken string) (service.Payload, error)
	VerifyRefreshTokenFunc   func(refreshToken string) (service.Payload, error)
	GenerateAccessTokenFunc  func(payload service.Payload) (string, error)
	GenerateRefreshTokenFunc func(payload service.Payload) (string, error)
}

func (m *MockAuthService) Register(ctx context.Context, req dto.SignUpRequest) (domain.User, error) {
	return m.RegisterFunc(ctx, req)
}

func (m *MockAuthService) Login(ctx context.Context, req dto.LoginRequest) (domain.User, string, string, error) {
	return m.LoginFunc(ctx, req)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	return m.RefreshTokenFunc(ctx, refreshToken)
}

func (m *MockAuthService) GetUserByAccessToken(ctx context.Context, accessToken string) (domain.User, error) {
	return m.GetUserByAccessTokenFunc(ctx, accessToken)
}

func (m *MockAuthService) VerifyAccessToken(accessToken string) (service.Payload, error) {
	return m.VerifyAccessTokenFunc(accessToken)
}

func (m *MockAuthService) VerifyRefreshToken(refreshToken string) (service.Payload, error) {
	return m.VerifyRefreshTokenFunc(refreshToken)
}

func (m *MockAuthService) GenerateAccessToken(payload service.Payload) (string, error) {
	return m.GenerateAccessTokenFunc(payload)
}

func (m *MockAuthService) GenerateRefreshToken(payload service.Payload) (string, error) {
	return m.GenerateRefreshTokenFunc(payload)
}

func generateTestToken(t *testing.T, userID uuid.UUID, email string, secret string) string {
	token := jwt.New()
	_ = token.Set(jwt.IssuedAtKey, time.Now().Unix())
	_ = token.Set(jwt.ExpirationKey, time.Now().Add(time.Minute*15).Unix())
	_ = token.Set("payload", service.Payload{
		UserID: userID,
		Email:  email,
	})

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, []byte(secret)))
	require.NoError(t, err)
	return string(signed)
}

func TestAuthHandler_SignUp(t *testing.T) {
	// Create a mock auth service.
	mockAuthService := &MockAuthService{
		RegisterFunc: func(ctx context.Context, req dto.SignUpRequest) (domain.User, error) {
			if req.Email == "existing@example.com" {
				return domain.User{}, fmt.Errorf("user with this email already exists")
			}
			return domain.User{
				ID:        uuid.New(),
				Email:     req.Email,
				FirstName: req.FirstName,
				LastName:  req.LastName,
				Role:      "user",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	// Create a new auth handler with the mock service.
	authHandler := NewAuthHandler(mockAuthService, service.NewValidatorService())

	// Create a new Fiber app for testing.
	app := fiber.New()

	// Setup the handler.
	app.Post("/signup", authHandler.SignUp)

	t.Run("Successful registration", func(t *testing.T) {
		// Create a new HTTP request.
		requestBody, _ := json.Marshal(dto.SignUpRequest{
			Email:     "test@example.com",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
		})
		req, _ := http.NewRequest("POST", "/signup", bytes.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Perform the request.
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("Malformed request body", func(t *testing.T) {
		// Create a new HTTP request.
		req, _ := http.NewRequest("POST", "/signup", bytes.NewReader([]byte("{")))
		req.Header.Set("Content-Type", "application/json")

		// Perform the request.
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Invalid request data", func(t *testing.T) {
		// Create a new HTTP request.
		requestBody, _ := json.Marshal(dto.SignUpRequest{
			Email:     "invalid-email",
			Password:  "short",
			FirstName: "",
			LastName:  "",
		})
		req, _ := http.NewRequest("POST", "/signup", bytes.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Perform the request.
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("User already exists", func(t *testing.T) {
		// Create a new HTTP request.
		requestBody, _ := json.Marshal(dto.SignUpRequest{
			Email:     "existing@example.com",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
		})
		req, _ := http.NewRequest("POST", "/signup", bytes.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Perform the request.
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestAuthHandler_Me(t *testing.T) {
	userID := uuid.New()
	testEmail := "test@example.com"
	testSecret := "test-secret-key-that-is-long-enough"

	// Create a mock auth service.
	mockAuthService := &MockAuthService{
		VerifyAccessTokenFunc: func(tokenStr string) (service.Payload, error) {
			token, err := jwt.Parse(
				[]byte(tokenStr),
				jwt.WithKey(jwa.HS256, []byte(testSecret)),
				jwt.WithValidate(true),
			)
			if err != nil {
				return service.Payload{}, fmt.Errorf("invalid token: %w", err)
			}

			payloadMap, ok := token.Get("payload")
			if !ok {
				return service.Payload{}, fmt.Errorf("missing payload")
			}

			data, err := json.Marshal(payloadMap)
			if err != nil {
				return service.Payload{}, fmt.Errorf("failed to marshal payload")
			}

			var payload service.Payload
			if err := json.Unmarshal(data, &payload); err != nil {
				return service.Payload{}, fmt.Errorf("failed to unmarshal payload")
			}

			return payload, nil
		},
		GetUserByAccessTokenFunc: func(ctx context.Context, accessToken string) (domain.User, error) {
			_, err := jwt.Parse(
				[]byte(accessToken),
				jwt.WithKey(jwa.HS256, []byte(testSecret)),
				jwt.WithValidate(true),
			)
			if err != nil {
				return domain.User{}, fmt.Errorf("invalid token: %w", err)
			}

			return domain.User{
				ID:        userID,
				Email:     testEmail,
				FirstName: "Test",
				LastName:  "User",
				Role:      "user",
			}, nil
		},
	}

	// Create a new auth handler with the mock service.
	authHandler := NewAuthHandler(mockAuthService, service.NewValidatorService())

	// Create a new Fiber app for testing.
	app := fiber.New()

	// Setup the middleware and handler.
	app.Get("/me", middleware.AuthMiddleware(mockAuthService), authHandler.Me)

	t.Run("Valid token", func(t *testing.T) {
		// Generate a real JWT for the test.
		validToken := generateTestToken(t, userID, testEmail, testSecret)

		// Create a new HTTP request.
		req, _ := http.NewRequest("GET", "/me", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)

		// Perform the request.
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Read the response body.
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		// Unmarshal the response body.
		var userResponse map[string]interface{}
		err = json.Unmarshal(body, &userResponse)
		require.NoError(t, err)

		// Assert the response body.
		assert.Equal(t, userID.String(), userResponse["id"])
		assert.Equal(t, testEmail, userResponse["email"])
	})

	t.Run("Malformed token", func(t *testing.T) {
		// Create a new HTTP request.
		req, _ := http.NewRequest("GET", "/me", nil)
		req.Header.Set("Authorization", "Bearer invalid-token-string")

		// Perform the request.
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Token with wrong secret", func(t *testing.T) {
		// Generate a token with a different secret.
		wrongToken := generateTestToken(t, userID, testEmail, "some-other-secret")

		// Create a new HTTP request.
		req, _ := http.NewRequest("GET", "/me", nil)
		req.Header.Set("Authorization", "Bearer "+wrongToken)

		// Perform the request.
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
