package utils

import "testing"

func TestValidatePassword(t *testing.T) {
	testCases := []struct {
		name          string
		password      string
		expectedError string
	}{
		{
			name:          "Valid password",
			password:      "ValidPass123",
			expectedError: "",
		},
		{
			name:          "Password too short",
			password:      "Vp1",
			expectedError: "password must be between 8 and 30 characters",
		},
		{
			name:          "Password too long",
			password:      "AVeryLongPasswordThatIsDefinitelyMoreThanThirtyCharacters1",
			expectedError: "password must be between 8 and 30 characters",
		},
		{
			name:          "No uppercase letter",
			password:      "validpass123",
			expectedError: "password must contain at least one uppercase and one lowercase letter",
		},
		{
			name:          "No lowercase letter",
			password:      "VALIDPASS123",
			expectedError: "password must contain at least one uppercase and one lowercase letter",
		},
		{
			name:          "No digit",
			password:      "ValidPassword",
			expectedError: "password must contain at least one digit",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePassword(tc.password)
			if tc.expectedError == "" {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error '%s', but got none", tc.expectedError)
				} else if err.Error() != tc.expectedError {
					t.Errorf("Expected error '%s', but got: %v", tc.expectedError, err)
				}
			}
		})
	}
}
