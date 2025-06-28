package repository

import (
	"errors"
	"fmt"
)

// Common repository errors
var (
	// Generic errors
	ErrNotFound      = errors.New("record not found")
	ErrAlreadyExists = errors.New("record already exists")
	ErrInvalidInput  = errors.New("invalid input provided")
	ErrDatabase      = errors.New("database operation failed")

	// User-specific errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")

	// Bank account errors
	ErrBankAccountNotFound      = errors.New("bank account not found")
	ErrBankAccountAlreadyExists = errors.New("bank account already exists")
	ErrInvalidAccountID         = errors.New("invalid account ID")

	// GoCardless item errors
	ErrGclItemNotFound      = errors.New("gocardless item not found")
	ErrGclItemAlreadyExists = errors.New("gocardless item already exists")
	ErrInvalidAccessToken   = errors.New("invalid access token")

	// Requisition errors
	ErrRequisitionNotFound      = errors.New("requisition not found")
	ErrRequisitionAlreadyExists = errors.New("requisition already exists")
	ErrInvalidRequisitionID     = errors.New("invalid requisition ID")

	// Transaction errors
	ErrTransactionNotFound    = errors.New("transaction not found")
	ErrInvalidTransactionData = errors.New("invalid transaction data")

	// Authorization errors
	ErrUnauthorized = errors.New("unauthorized access")
	ErrForbidden    = errors.New("forbidden operation")
)

// RepositoryError wraps repository errors with additional context
type RepositoryError struct {
	Operation string
	Entity    string
	Err       error
	Context   map[string]interface{}
}

func (e *RepositoryError) Error() string {
	if len(e.Context) > 0 {
		return fmt.Sprintf("repository error: %s %s failed: %v (context: %v)",
			e.Operation, e.Entity, e.Err, e.Context)
	}
	return fmt.Sprintf("repository error: %s %s failed: %v",
		e.Operation, e.Entity, e.Err)
}

func (e *RepositoryError) Unwrap() error {
	return e.Err
}

// NewRepositoryError creates a new repository error with context
func NewRepositoryError(operation, entity string, err error, context ...map[string]interface{}) *RepositoryError {
	var ctx map[string]interface{}
	if len(context) > 0 {
		ctx = context[0]
	}

	return &RepositoryError{
		Operation: operation,
		Entity:    entity,
		Err:       err,
		Context:   ctx,
	}
}

// Helper functions to create specific repository errors
func NewUserError(operation string, err error, context ...map[string]interface{}) *RepositoryError {
	return NewRepositoryError(operation, "user", err, context...)
}

func NewBankAccountError(operation string, err error, context ...map[string]interface{}) *RepositoryError {
	return NewRepositoryError(operation, "bank_account", err, context...)
}

func NewGclItemError(operation string, err error, context ...map[string]interface{}) *RepositoryError {
	return NewRepositoryError(operation, "gocardless_item", err, context...)
}

func NewRequisitionError(operation string, err error, context ...map[string]interface{}) *RepositoryError {
	return NewRepositoryError(operation, "requisition", err, context...)
}

func NewTransactionError(operation string, err error, context ...map[string]interface{}) *RepositoryError {
	return NewRepositoryError(operation, "transaction", err, context...)
}

// IsNotFoundError checks if an error is a "not found" error
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	// Check direct error
	if errors.Is(err, ErrNotFound) ||
		errors.Is(err, ErrUserNotFound) ||
		errors.Is(err, ErrBankAccountNotFound) ||
		errors.Is(err, ErrGclItemNotFound) ||
		errors.Is(err, ErrRequisitionNotFound) ||
		errors.Is(err, ErrTransactionNotFound) {
		return true
	}

	// Check wrapped repository error
	var repoErr *RepositoryError
	if errors.As(err, &repoErr) {
		return IsNotFoundError(repoErr.Err)
	}

	return false
}

// IsAlreadyExistsError checks if an error is an "already exists" error
func IsAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}

	// Check direct error
	if errors.Is(err, ErrAlreadyExists) ||
		errors.Is(err, ErrUserAlreadyExists) ||
		errors.Is(err, ErrBankAccountAlreadyExists) ||
		errors.Is(err, ErrGclItemAlreadyExists) ||
		errors.Is(err, ErrRequisitionAlreadyExists) {
		return true
	}

	// Check wrapped repository error
	var repoErr *RepositoryError
	if errors.As(err, &repoErr) {
		return IsAlreadyExistsError(repoErr.Err)
	}

	return false
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrInvalidInput) ||
		errors.Is(err, ErrInvalidCredentials) ||
		errors.Is(err, ErrInvalidAccountID) ||
		errors.Is(err, ErrInvalidAccessToken) ||
		errors.Is(err, ErrInvalidRequisitionID) ||
		errors.Is(err, ErrInvalidTransactionData)
}

// IsAuthorizationError checks if an error is an authorization error
func IsAuthorizationError(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrForbidden)
}
