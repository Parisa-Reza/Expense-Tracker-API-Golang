package controllers

import (
	"net/http"
	"strconv"
	"testing"
)

func TestCreateExpense(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "valid expense",
			input:      `{"title":"Lunch","amount":350.50,"category":"Food","expense_date":"2025-06-10"}`,
			wantStatus: http.StatusCreated,
			wantMsg:    "Expense created successfully",
		},
		{
			name:       "missing title",
			input:      `{"amount":350.50,"category":"Food","expense_date":"2025-06-10"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Title, category, and expense_date are required",
		},
		{
			name:       "invalid category",
			input:      `{"title":"Lunch","amount":350.50,"category":"InvalidCat","expense_date":"2025-06-10"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Invalid expense category",
		},
		{
			name:       "invalid amount",
			input:      `{"title":"Lunch","amount":0,"category":"Food","expense_date":"2025-06-10"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Amount must be greater than zero",
		},
		{
			name:       "invalid date",
			input:      `{"title":"Lunch","amount":350.50,"category":"Food","expense_date":"06-10-2025"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "expense_date must use YYYY-MM-DD format",
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
			user := createTestUser(t)

			response := performRequest(http.MethodPost, "/api/v1/expenses", tt.input, strconv.Itoa(user.ID))
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
