package controllers

import (
	"net/http"
	"strconv"
	"testing"
)

func TestUpdateExpense(t *testing.T) {
	tests := []struct {
		name       string
		target     string
		seed       bool
		input      string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "valid update",
			seed:       true,
			input:      `{"title":"Dinner","amount":500,"category":"Food","expense_date":"2025-06-11"}`,
			wantStatus: http.StatusOK,
			wantMsg:    "Expense updated successfully",
		},
		{
			name:       "invalid body",
			seed:       true,
			input:      `{bad json`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Invalid request body",
		},
		{
			name:       "missing expense",
			target:     "/api/v1/expenses/99",
			input:      `{"title":"Dinner","amount":500,"category":"Food","expense_date":"2025-06-11"}`,
			wantStatus: http.StatusNotFound,
			wantMsg:    "Expense not found",
		},
		{
			name:       "invalid id",
			target:     "/api/v1/expenses/abc",
			input:      `{"title":"Dinner","amount":500,"category":"Food","expense_date":"2025-06-11"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Invalid expense ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupControllerTest(t)
			user := createTestUser(t)
			target := tt.target
			if tt.seed {
				expense := createTestExpense(t, user.ID)
				target = "/api/v1/expenses/" + strconv.Itoa(expense.ID)
			}

			response := performRequest(http.MethodPut, target, tt.input, strconv.Itoa(user.ID))
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
