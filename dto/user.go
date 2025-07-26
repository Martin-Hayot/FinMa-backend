package dto

import "github.com/google/uuid"

// UserResponse represents user data for API responses
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Role      string    `json:"role"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
}

// UpdateProfileRequest represents data needed to update a user profile
type UpdateProfileRequest struct {
	FirstName string `json:"firstName" validate:"omitempty,min=2,max=50"`
	LastName  string `json:"lastName" validate:"omitempty,min=2,max=50"`
	Email     string `json:"email" validate:"omitempty,email"`
}

// ChangePasswordRequest represents data needed to change a user's password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=NewPassword"`
}

// UserPreferencesResponse represents a user's preferences
type UserPreferencesResponse struct {
	Currency             string `json:"currency"`
	Theme                string `json:"theme"`
	NotificationsEnabled bool   `json:"notificationsEnabled"`
	BudgetAlerts         bool   `json:"budgetAlerts"`
	WeeklyReports        bool   `json:"weeklyReports"`
}

// UpdatePreferencesRequest represents data needed to update a user's preferences
type UpdatePreferencesRequest struct {
	Currency             string `json:"currency" validate:"omitempty,len=3"`
	Theme                string `json:"theme" validate:"omitempty,oneof=light dark system"`
	NotificationsEnabled *bool  `json:"notificationsEnabled"`
	BudgetAlerts         *bool  `json:"budgetAlerts"`
	WeeklyReports        *bool  `json:"weeklyReports"`
}

type UserSummaryResponse struct {
	TotalIncome            float64 `json:"totalIncome"`
	TotalExpenses          float64 `json:"totalExpenses"`
	TotalSavings           float64 `json:"totalSavings"`
	AverageMonthlyIncome   float64 `json:"averageMonthlyIncome"`
	AverageMonthlyExpenses float64 `json:"averageMonthlyExpenses"`
}
