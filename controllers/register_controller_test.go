package controllers

import (
	"net/http"
	"testing"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name       string
		seedUser   bool
		input      string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "valid user",
			input:      `{"name":"Parisa","email":"parisa@example.com","password":"secret1"}`,
			wantStatus: http.StatusCreated,
			wantMsg:    "User registered successfully",
		},
		{
			name:       "invalid json",
			input:      `{bad json`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Invalid request body",
		},
		{
			name:       "missing name",
			input:      `{"email":"parisa@example.com","password":"secret1"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Name is required",
		},
		{
			name:       "missing email",
			input:      `{"name":"Parisa","password":"secret1"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Email is required",
		},
		{
			name:       "invalid email",
			input:      `{"name":"Parisa","email":"not-email","password":"secret1"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Email must be valid",
		},
		{
			name:       "missing password",
			input:      `{"name":"Parisa","email":"parisa@example.com"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Password is required",
		},
		{
			name:       "short password",
			input:      `{"name":"Parisa","email":"parisa@example.com","password":"123"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Password must be at least 6 characters",
		},
		{
			name:       "duplicate email",
			seedUser:   true,
			input:      `{"name":"Other","email":"test@example.com","password":"secret1"}`,
			wantStatus: http.StatusConflict,
			wantMsg:    "Email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupControllerTest(t)
			if tt.seedUser {
				createTestUser(t)
			}

			response := performRequest(http.MethodPost, "/api/v1/auth/register", tt.input, "")
			got := decodeAPIResponse(t, response)

			if response.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", response.Code, tt.wantStatus, response.Body.String())
			}
			if got.Message != tt.wantMsg {
				t.Fatalf("message = %q, want %q", got.Message, tt.wantMsg)
			}
		})
	}
}
