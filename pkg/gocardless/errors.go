package gocardless

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GoCardlessError represents an error response from the GoCardless API
type GoCardlessError struct {
	StatusCode int                    `json:"status_code"`
	Reference  map[string]interface{} `json:"reference,omitempty"`
	Detail     string                 `json:"detail,omitempty"`
	Summary    string                 `json:"summary,omitempty"`
	Type       string                 `json:"type,omitempty"`
}

func (e *GoCardlessError) Error() string {
	if e.Reference != nil {
		if summary, ok := e.Reference["summary"].(string); ok {
			if detail, ok := e.Reference["detail"].(string); ok {
				return fmt.Sprintf("GoCardless API error (status %d): %s - %s", e.StatusCode, summary, detail)
			}
			return fmt.Sprintf("GoCardless API error (status %d): %s", e.StatusCode, summary)
		}
	}

	if e.Summary != "" {
		if e.Detail != "" {
			return fmt.Sprintf("GoCardless API error (status %d): %s - %s", e.StatusCode, e.Summary, e.Detail)
		}
		return fmt.Sprintf("GoCardless API error (status %d): %s", e.StatusCode, e.Summary)
	}

	return fmt.Sprintf("GoCardless API error (status %d)", e.StatusCode)
}

// IsConflictError checks if the error is a conflict (409) or reference already exists (400)
func (e *GoCardlessError) IsConflictError() bool {
	if e.StatusCode == http.StatusConflict {
		return true
	}

	if e.StatusCode == http.StatusBadRequest && e.Reference != nil {
		if summary, ok := e.Reference["summary"].(string); ok {
			return summary == "Client reference must be unique"
		}
	}

	return false
}

// parseGoCardlessError parses an error response from GoCardless API
func parseGoCardlessError(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read error response body: %w", err)
	}

	var gcError GoCardlessError
	if err := json.Unmarshal(body, &gcError); err != nil {
		// If we can't parse the error response, return a generic error
		return fmt.Errorf("GoCardless API error (status %d): %s", resp.StatusCode, string(body))
	}

	gcError.StatusCode = resp.StatusCode
	return &gcError
}
