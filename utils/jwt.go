package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// GenerateAccessToken generates a new JWT access token.
// The payload is the data that will be stored in the token.
// The function returns the signed token as a string.
func GenerateAccessToken[T any](payload T) (string, error) {
	// Retrieve the secret key from environment variables
	secretKey := os.Getenv("ACCESS_TOKEN_SECRET")
	if secretKey == "" {
		return "", fmt.Errorf("access token secret is not set")
	}

	// Generate a new JWT token
	token := jwt.New()

	// Set the token claims
	token.Set("payload", payload)
	token.Set(jwt.IssuedAtKey, time.Now().Unix())
	token.Set(jwt.ExpirationKey, time.Now().Add(time.Hour*1).Unix())
	token.Set(jwt.IssuerKey, "FinMa")
	token.Set(jwt.SubjectKey, "access")
	token.Set(jwt.AudienceKey, "users")

	// Sign the token
	signedToken, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, []byte(secretKey)))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return string(signedToken), nil
}

func GenerateRefreshToken() {
	// Generate a new JWT token
}

// VerifyAccessToken verifies the JWT access token.
// The function returns the payload stored in the token.
func VerifyAccessToken() {
	// Verify the JWT token
}

func VerifyRefreshToken() {
	// Verify the JWT token
}
