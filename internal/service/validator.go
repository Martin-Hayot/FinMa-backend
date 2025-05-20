package service

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidatorService defines operations for validating data
type ValidatorService interface {
	Validate(data interface{}) error
	RegisterCustomValidation(tag string, fn validator.Func) error
}

type validatorService struct {
	validate *validator.Validate
}

// NewValidatorService creates a new validator service
func NewValidatorService() ValidatorService {
	validate := validator.New()

	// Register validation for structs to use json tags as field names in errors
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validations
	registerCustomValidations(validate)

	return &validatorService{
		validate: validate,
	}
}

// Validate validates the provided data against its validation tags
func (s *validatorService) Validate(data interface{}) error {
	if data == nil {
		return nil
	}

	err := s.validate.Struct(data)
	if err != nil {
		// Convert validation errors to more user-friendly format
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := make([]string, 0, len(validationErrors))

			for _, e := range validationErrors {
				errors = append(errors, formatValidationError(e))
			}

			return fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
		}
		return err
	}

	return nil
}

// RegisterCustomValidation registers a custom validation function
func (s *validatorService) RegisterCustomValidation(tag string, fn validator.Func) error {
	return s.validate.RegisterValidation(tag, fn)
}

// Helper function to format validation errors in a user-friendly way
func formatValidationError(err validator.FieldError) string {
	field := err.Field()

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, err.Param())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s characters", field, err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, err.Param())
	case "eqfield":
		return fmt.Sprintf("%s must be equal to %s", field, err.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, err.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, err.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, err.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, err.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "datetime":
		return fmt.Sprintf("%s must be a valid date in format %s", field, err.Param())
	default:
		return fmt.Sprintf("%s failed validation: %s", field, err.Tag())
	}
}

// Register custom validations for specific FinMa business rules
func registerCustomValidations(validate *validator.Validate) {
	// Example: Custom validation for transaction categories
	validate.RegisterValidation("valid_transaction_category", validateTransactionCategory)

	// Example: Custom validation for currency codes
	validate.RegisterValidation("currency", validateCurrencyCode)

	// Example: Custom validation for password strength
	validate.RegisterValidation("strong_password", validateStrongPassword)
}

// Custom validation for transaction categories
func validateTransactionCategory(fl validator.FieldLevel) bool {
	category := fl.Field().String()
	validCategories := []string{
		"food", "transport", "housing", "utilities",
		"entertainment", "health", "education", "shopping",
		"personal", "debt", "savings", "income", "other",
	}

	for _, validCategory := range validCategories {
		if category == validCategory {
			return true
		}
	}
	return false
}

// Custom validation for currency codes
func validateCurrencyCode(fl validator.FieldLevel) bool {
	code := fl.Field().String()
	if len(code) != 3 {
		return false
	}

	// Check if it's all uppercase letters
	for _, r := range code {
		if r < 'A' || r > 'Z' {
			return false
		}
	}

	// Here you could also check against a list of valid currency codes
	// For simplicity, we're just checking format
	return true
}

// Custom validation for strong passwords
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check length
	if len(password) < 8 {
		return false
	}

	// Check for at least one uppercase letter
	hasUpper := false
	for _, c := range password {
		if c >= 'A' && c <= 'Z' {
			hasUpper = true
			break
		}
	}

	// Check for at least one lowercase letter
	hasLower := false
	for _, c := range password {
		if c >= 'a' && c <= 'z' {
			hasLower = true
			break
		}
	}

	// Check for at least one digit
	hasDigit := false
	for _, c := range password {
		if c >= '0' && c <= '9' {
			hasDigit = true
			break
		}
	}

	// Check for at least one special character
	hasSpecial := false
	specialChars := "!@#$%^&*()-_=+[]{}|;:'\",.<>/?"
	for _, c := range password {
		if strings.ContainsRune(specialChars, c) {
			hasSpecial = true
			break
		}
	}

	// Password is strong if it has all requirements
	return hasUpper && hasLower && hasDigit && hasSpecial
}
