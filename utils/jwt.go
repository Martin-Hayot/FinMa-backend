package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var (
	// AccessTokenSecret is the secret key used to sign the access token.
	AccessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
	// RefreshTokenSecret is the secret key used to sign the refresh token.
	RefreshTokenSecret = os.Getenv("REFRESH_TOKEN_SECRET")
)

// GenerateAccessToken generates a new JWT access token.
// The payload is the data that will be stored in the token.
// The function returns the signed token as a string.
func GenerateAccessToken[T any](payload T) (string, error) {
	if AccessTokenSecret == "" {
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
	signedToken, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, []byte(AccessTokenSecret)))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return string(signedToken), nil
}

func GenerateRefreshToken[T any](payload T) (string, error) {
	if RefreshTokenSecret == "" {
		return "", fmt.Errorf("refresh token secret is not set")
	}

	// Generate a new JWT token
	token := jwt.New()

	// Set the token claims
	token.Set("payload", payload)
	token.Set(jwt.IssuedAtKey, time.Now().Unix())
	token.Set(jwt.ExpirationKey, time.Now().Add(time.Hour*24).Unix())
	token.Set(jwt.IssuerKey, "FinMa")
	token.Set(jwt.SubjectKey, "refresh")
	token.Set(jwt.AudienceKey, "users")

	// Sign the token
	signedToken, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, []byte(RefreshTokenSecret)))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return string(signedToken), nil
}

// VerifyAccessToken verifies the JWT access token.
// The function returns the payload stored in the token.
func VerifyAccessToken(tokenString string) (interface{}, error) {
	if AccessTokenSecret == "" {
		return nil, fmt.Errorf("access token secret is not set")
	}

	// Parse the token
	token, err := jwt.Parse([]byte(tokenString), jwt.WithKey(jwa.HS256, []byte(AccessTokenSecret)), jwt.WithValidate(true))

	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	// Retrieve the payload from the token

	payload, ok := token.Get("payload")

	if !ok {
		return nil, fmt.Errorf("failed to retrieve payload from token")
	}

	fmt.Println(payload)

	return payload, nil
}

// VerifyRefreshToken verifies the JWT refresh token.
// The function returns the payload stored in the token.
func VerifyRefreshToken(tokenString string) (interface{}, error) {
	if RefreshTokenSecret == "" {
		return "", fmt.Errorf("refresh token secret is not set")
	}

	// Parse the token
	token, err := jwt.Parse([]byte(tokenString), jwt.WithKey(jwa.HS256, []byte(RefreshTokenSecret)), jwt.WithValidate(true))

	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	// Retrieve the payload from the token

	payload, ok := token.Get("payload")

	if !ok {
		return nil, fmt.Errorf("failed to retrieve payload from token")
	}

	fmt.Println(payload)

	return payload, nil
}
