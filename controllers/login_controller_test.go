package controllers

import (
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name       string
		seedUser   bool
		input      string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "valid login",
			seedUser:   true,
			input:      `{"email":"test@example.com","password":"secret1"}`,
			wantStatus: http.StatusOK,
			wantMsg:    "Login successful",
		},
		{
			name:       "wrong password",
			seedUser:   true,
			input:      `{"email":"test@example.com","password":"badpass"}`,
			wantStatus: http.StatusUnauthorized,
			wantMsg:    "Invalid email or password",
		},
		{
			name:       "unknown email",
			input:      `{"email":"missing@example.com","password":"secret1"}`,
			wantStatus: http.StatusUnauthorized,
			wantMsg:    "Invalid email or password",
		},
		{
			name:       "invalid json",
			input:      `{bad json`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupControllerTest(t)
			if tt.seedUser {
				createTestUser(t)
			}

			response := performRequest(http.MethodPost, "/api/v1/auth/login", tt.input, "")
			got := decodeAPIResponse(t, response)

			if response.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, tt.wantStatus)
			}
			if got.Message != tt.wantMsg {
				t.Fatalf("message = %q, want %q", got.Message, tt.wantMsg)
			}
		})
	}
}
