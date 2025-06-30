package gocardless

import (
	"FinMa/dto"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	BaseURL              = "https://bankaccountdata.gocardless.com/api/v2"
	TokenEndpoint        = "/token/new/"
	RefreshEndpoint      = "/token/refresh/"
	InstitutionsEndpoint = "/institutions/"
	RequisitionsEndpoint = "/requisitions/"
	AccountsEndpoint     = "/accounts/"
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

func (c *Client) CreateRequisition(ctx context.Context, UserID uuid.UUID, institutionID, redirectURL string) (*dto.GoCardlessCreateRequisitionResponse, error) {
	accessToken, err := c.GetValidAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, RequisitionsEndpoint)

	// Create a unique reference by combining UserID and institutionID
	uniqueReference := fmt.Sprintf("%s_%s", UserID.String(), institutionID)

	jsonData, err := json.Marshal(dto.GoCardlessCreateRequisitionRequest{
		InstitutionID: institutionID,
		RedirectURL:   redirectURL,
		Reference:     uniqueReference, // Now unique per user-institution combination
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, parseGoCardlessError(resp)
	}

	var linkResp dto.GoCardlessCreateRequisitionResponse
	if err := json.NewDecoder(resp.Body).Decode(&linkResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &linkResp, nil
}

func (c *Client) GetRequisition(ctx context.Context, userID uuid.UUID, requisitionID string) (*dto.GoCardlessGetRequisitionResponse, error) {
	accessToken, err := c.GetValidAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s%s%s/", c.BaseURL, RequisitionsEndpoint, requisitionID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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
		return nil, parseGoCardlessError(resp)
	}

	var requisition dto.GoCardlessGetRequisitionResponse
	if err := json.NewDecoder(resp.Body).Decode(&requisition); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Security check: Verify the requisition belongs to the user
	// The reference should start with the userID (format: "userID_institutionID")
	expectedPrefix := userID.String() + "_"
	if !strings.HasPrefix(requisition.Reference, expectedPrefix) {
		return nil, fmt.Errorf("access denied: requisition does not belong to user")
	}

	return &requisition, nil
}

// GetValidAccessToken returns a valid access token, refreshing if necessary
func (c *Client) GetValidAccessToken(ctx context.Context) (string, error) {
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
		if err := c.refreshAccessTokenInternal(ctx); err == nil {
			return c.AccessToken, nil
		}
		// Refresh failed, fall through to get new token
	}

	// Get new token
	return c.getNewAccessTokenInternal(ctx)
}

// Internal method to get new access token (must be called with write lock)
func (c *Client) getNewAccessTokenInternal(ctx context.Context) (string, error) {
	tokenResp, err := c.getAccessTokenRequest(ctx)
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
func (c *Client) refreshAccessTokenInternal(ctx context.Context) error {
	tokenResp, err := c.refreshAccessTokenRequest(ctx, c.RefreshToken)
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
func (c *Client) getAccessTokenRequest(ctx context.Context) (*dto.TokenResponse, error) {
	jsonData, err := json.Marshal(dto.TokenRequest{
		SecretID:  c.SecretID,
		SecretKey: c.SecretKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal token request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+TokenEndpoint, bytes.NewBuffer(jsonData))
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

func (c *Client) refreshAccessTokenRequest(ctx context.Context, refreshToken string) (*dto.TokenResponse, error) {
	refreshReq := dto.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	jsonData, err := json.Marshal(refreshReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal refresh request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+RefreshEndpoint, bytes.NewBuffer(jsonData))
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

func (c *Client) GetInstitutions(ctx context.Context, countryCode string) ([]dto.Institution, error) {
	accessToken, err := c.GetValidAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s%s?country=%s", c.BaseURL, InstitutionsEndpoint, countryCode)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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

// GetAccountDetails retrieves the details of a specific bank account
func (c *Client) GetAccountDetails(ctx context.Context, accountID string) (*dto.AccountDetails, error) {
	accessToken, err := c.GetValidAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s%s%s/details/", c.BaseURL, AccountsEndpoint, accountID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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
		return nil, parseGoCardlessError(resp)
	}

	var details dto.AccountDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &details, nil
}

// GetAccountBalances retrieves the balances of a specific bank account
func (c *Client) GetAccountBalances(ctx context.Context, accountID string) (*dto.AccountBalances, error) {
	accessToken, err := c.GetValidAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s%s%s/balances/", c.BaseURL, AccountsEndpoint, accountID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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
		return nil, parseGoCardlessError(resp)
	}

	var balances dto.AccountBalances
	if err := json.NewDecoder(resp.Body).Decode(&balances); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &balances, nil
}

// GetAccountTransactions retrieves the transactions of a specific bank account
func (c *Client) GetAccountTransactions(ctx context.Context, accountID string) (*dto.AccountTransactions, error) {
	accessToken, err := c.GetValidAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s%s%s/transactions/", c.BaseURL, AccountsEndpoint, accountID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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
		return nil, parseGoCardlessError(resp)
	}

	var transactions dto.AccountTransactions
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &transactions, nil
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
func (c *Client) GetAccessToken(ctx context.Context) (*dto.TokenResponse, error) {
	return c.getAccessTokenRequest(ctx)
}

func (c *Client) RefreshAccessToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error) {
	return c.refreshAccessTokenRequest(ctx, refreshToken)
}
