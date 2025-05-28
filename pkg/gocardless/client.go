package gocardless

import (
	"FinMa/dto"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	BaseURL              = "https://bankaccountdata.gocardless.com/api/v2"
	TokenEndpoint        = "/token/new/"
	RefreshEndpoint      = "/token/refresh/"
	InstitutionsEndpoint = "/institutions/"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	SecretID   string
	SecretKey  string

	// Token management with mutex for thread safety
	mu             sync.RWMutex
	AccessToken    string
	RefreshToken   string
	AccessExpires  time.Time
	RefreshExpires time.Time
}

func NewClient(secretID, secretKey string) *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		BaseURL:   BaseURL,
		SecretID:  secretID,
		SecretKey: secretKey,
	}
}

// GetValidAccessToken returns a valid access token, refreshing if necessary
func (c *Client) GetValidAccessToken() (string, error) {
	c.mu.RLock()
	// Check if we have a valid access token (with 5 minute buffer)
	if c.AccessToken != "" && time.Now().Add(5*time.Minute).Before(c.AccessExpires) {
		token := c.AccessToken
		c.mu.RUnlock()
		return token, nil
	}
	c.mu.RUnlock()

	// Need to refresh or get new token
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if c.AccessToken != "" && time.Now().Add(5*time.Minute).Before(c.AccessExpires) {
		return c.AccessToken, nil
	}

	// Try to refresh if we have a valid refresh token
	if c.RefreshToken != "" && time.Now().Before(c.RefreshExpires) {
		if err := c.refreshAccessTokenInternal(); err == nil {
			return c.AccessToken, nil
		}
		// Refresh failed, fall through to get new token
	}

	// Get new token
	return c.getNewAccessTokenInternal()
}

// Internal method to get new access token (must be called with write lock)
func (c *Client) getNewAccessTokenInternal() (string, error) {
	tokenResp, err := c.getAccessTokenRequest()
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}

	c.AccessToken = tokenResp.Access
	c.RefreshToken = tokenResp.Refresh
	c.AccessExpires = time.Now().Add(time.Duration(tokenResp.AccessExpires) * time.Second)
	c.RefreshExpires = time.Now().Add(time.Duration(tokenResp.RefreshExpires) * time.Second)

	return c.AccessToken, nil
}

// Internal method to refresh access token (must be called with write lock)
func (c *Client) refreshAccessTokenInternal() error {
	tokenResp, err := c.refreshAccessTokenRequest(c.RefreshToken)
	if err != nil {
		return fmt.Errorf("failed to refresh access token: %w", err)
	}

	c.AccessToken = tokenResp.Access
	c.RefreshToken = tokenResp.Refresh
	c.AccessExpires = time.Now().Add(time.Duration(tokenResp.AccessExpires) * time.Second)
	c.RefreshExpires = time.Now().Add(time.Duration(tokenResp.RefreshExpires) * time.Second)

	return nil
}

// Raw token request methods
func (c *Client) getAccessTokenRequest() (*dto.TokenResponse, error) {
	jsonData, err := json.Marshal(dto.TokenRequest{
		SecretID:  c.SecretID,
		SecretKey: c.SecretKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal token request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+TokenEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var tokenResp dto.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &tokenResp, nil
}

func (c *Client) refreshAccessTokenRequest(refreshToken string) (*dto.TokenResponse, error) {
	refreshReq := dto.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	jsonData, err := json.Marshal(refreshReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal refresh request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+RefreshEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var tokenResp dto.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &tokenResp, nil
}

func (c *Client) GetInstitutions(countryCode string) ([]dto.Institution, error) {
	accessToken, err := c.GetValidAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s%s?country=%s", c.BaseURL, InstitutionsEndpoint, countryCode)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var institutions []dto.Institution
	if err := json.NewDecoder(resp.Body).Decode(&institutions); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return institutions, nil
}

// GetTokenStatus returns the current token status
func (c *Client) GetTokenStatus() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.AccessToken == "" {
		return map[string]interface{}{
			"has_token": false,
		}
	}

	now := time.Now()
	return map[string]interface{}{
		"has_token":           true,
		"access_token_valid":  now.Before(c.AccessExpires),
		"refresh_token_valid": now.Before(c.RefreshExpires),
		"access_expires_in":   int(c.AccessExpires.Sub(now).Seconds()),
		"refresh_expires_in":  int(c.RefreshExpires.Sub(now).Seconds()),
	}
}

// ClearTokens clears stored tokens
func (c *Client) ClearTokens() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.AccessToken = ""
	c.RefreshToken = ""
	c.AccessExpires = time.Time{}
	c.RefreshExpires = time.Time{}
}

// Legacy methods for backward compatibility
func (c *Client) GetAccessToken() (*dto.TokenResponse, error) {
	return c.getAccessTokenRequest()
}

func (c *Client) RefreshAccessToken(refreshToken string) (*dto.TokenResponse, error) {
	return c.refreshAccessTokenRequest(refreshToken)
}
