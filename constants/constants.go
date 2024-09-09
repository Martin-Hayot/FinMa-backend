package constants

// Constants for the application.

var TRANSACTION_TYPES = []string{"income", "expense"}

var TRANSACTION_CATEGORIES = []string{"food", "transport", "shopping", "bills", "others"}

var USER_ROLES = []string{"user", "admin"}

func GetTransactionTypes() []string {
	return append([]string(nil), TRANSACTION_TYPES...)
}

func GetTransactionCategories() []string {
	return append([]string(nil), TRANSACTION_CATEGORIES...)
}

func GetUserRoles() []string {
	return append([]string(nil), USER_ROLES...)
}
