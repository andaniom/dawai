package services

import (
	"context"
	"testing"
)

func TestGenerateTokenValidation(t *testing.T) {
	authService := &AuthService{}
	ctx := context.Background()

	tests := []struct {
		name      string
		userID    string
		schoolID  string
		roles     []string
		shouldErr bool
		errMsg    string
	}{
		{
			name:      "valid token with all claims",
			userID:    "user-123",
			schoolID:  "school-456",
			roles:     []string{"teacher"},
			shouldErr: false,
		},
		{
			name:      "missing userID",
			userID:    "",
			schoolID:  "school-456",
			roles:     []string{"teacher"},
			shouldErr: true,
			errMsg:    "missing or empty userID",
		},
		{
			name:      "missing schoolID",
			userID:    "user-123",
			schoolID:  "",
			roles:     []string{"teacher"},
			shouldErr: true,
			errMsg:    "missing or empty schoolID",
		},
		{
			name:      "both userID and schoolID missing",
			userID:    "",
			schoolID:  "",
			roles:     []string{},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := authService.GenerateToken(ctx, tt.userID, tt.schoolID, tt.roles)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Logf("Got error: %v", err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if token == "" {
					t.Errorf("Expected non-empty token")
				}
			}
		})
	}
}
