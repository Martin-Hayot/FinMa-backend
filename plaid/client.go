package plaid

import (
	"FinMa/config"
	"context"
	"fmt"
	"time"

	plaid "github.com/plaid/plaid-go/v31/plaid"
)

// Client wraps the Plaid API client with additional state
type Client struct {
	APIClient   *plaid.APIClient
	AccessToken string
	ItemID      string
}

// NewClient creates a new Plaid client
func NewClient(cfg *config.Config) *Client {
	// Create Plaid client configuration
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", cfg.PlaidClientID)
	configuration.AddDefaultHeader("PLAID-SECRET", cfg.PlaidSecret)
	configuration.UseEnvironment(cfg.PlaidEnv)

	return &Client{
		APIClient: plaid.NewAPIClient(configuration),
	}
}

// ExchangePublicToken exchanges a public token for an access token
func (c *Client) ExchangePublicToken(ctx context.Context, publicToken string) (string, string, error) {
	// Exchange the public_token for an access_token
	exchangeReq := plaid.NewItemPublicTokenExchangeRequest(publicToken)
	exchangeResp, _, err := c.APIClient.PlaidApi.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(*exchangeReq).Execute()
	if err != nil {
		return "", "", err
	}

	c.AccessToken = exchangeResp.GetAccessToken()
	c.ItemID = exchangeResp.GetItemId()

	return c.AccessToken, c.ItemID, nil
}

// CreateLinkToken creates a link token for initializing Plaid Link
func (c *Client) CreateLinkToken(ctx context.Context, cfg *config.Config) (string, error) {
	// Convert country codes to Plaid format
	countryCodes := []plaid.CountryCode{}
	for _, code := range cfg.PlaidCountryCodes {
		countryCodes = append(countryCodes, plaid.CountryCode(code))
	}

	// Convert products to Plaid format
	products := []plaid.Products{}
	for _, product := range cfg.PlaidProducts {
		products = append(products, plaid.Products(product))
	}

	// Create a unique client user ID (typically this would be your user's ID)
	clientUserID := fmt.Sprintf("user-%d", time.Now().Unix())

	// Create the link token request
	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: clientUserID,
	}

	request := plaid.NewLinkTokenCreateRequest(
		"Plaid Fiber App",
		"en",
		countryCodes,
		user,
	)
	request.SetProducts(products)

	if cfg.PlaidRedirectURI != "" {
		request.SetRedirectUri(cfg.PlaidRedirectURI)
	}

	// Execute the request
	resp, _, err := c.APIClient.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		return "", err
	}

	return resp.GetLinkToken(), nil
}

// GetAccounts retrieves accounts for the current access token
func (c *Client) GetAccounts(ctx context.Context) (*plaid.AccountsGetResponse, error) {
	request := plaid.NewAccountsGetRequest(c.AccessToken)
	response, _, err := c.APIClient.PlaidApi.AccountsGet(ctx).AccountsGetRequest(*request).Execute()
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetTransactions retrieves transactions for the current access token
func (c *Client) GetTransactions(ctx context.Context) ([]plaid.Transaction, error) {
	// Set cursor to empty to receive all historical updates
	var cursor *string
	var transactions []plaid.Transaction
	hasMore := true

	// Iterate through each page of transaction updates
	for hasMore {
		request := plaid.NewTransactionsSyncRequest(c.AccessToken)
		if cursor != nil {
			request.SetCursor(*cursor)
		}

		resp, _, err := c.APIClient.PlaidApi.TransactionsSync(ctx).TransactionsSyncRequest(*request).Execute()
		if err != nil {
			return nil, err
		}

		// Update cursor for the next request
		nextCursor := resp.GetNextCursor()
		cursor = &nextCursor

		// If no transactions are available yet, pause and try again
		if *cursor == "" {
			time.Sleep(1 * time.Second)
			continue
		}

		// Add this page of results
		transactions = append(transactions, resp.GetAdded()...)
		hasMore = resp.GetHasMore()
	}

	return transactions, nil
}

// GetItem retrieves item information and institution for the current access token
func (c *Client) GetItem(ctx context.Context, countryCodes []string) (*plaid.ItemGetResponse, *plaid.InstitutionsGetByIdResponse, error) {
	// Get item
	itemRequest := plaid.NewItemGetRequest(c.AccessToken)
	itemResp, _, err := c.APIClient.PlaidApi.ItemGet(ctx).ItemGetRequest(*itemRequest).Execute()
	if err != nil {
		return nil, nil, err
	}

	// Convert country codes
	plaidCountryCodes := []plaid.CountryCode{}
	for _, code := range countryCodes {
		plaidCountryCodes = append(plaidCountryCodes, plaid.CountryCode(code))
	}

	// Get institution
	institutionID := *itemResp.GetItem().InstitutionId.Get()
	institutionRequest := plaid.NewInstitutionsGetByIdRequest(institutionID, plaidCountryCodes)
	institutionResp, _, err := c.APIClient.PlaidApi.InstitutionsGetById(ctx).InstitutionsGetByIdRequest(*institutionRequest).Execute()
	if err != nil {
		return &itemResp, nil, err
	}

	return &itemResp, &institutionResp, nil
}
