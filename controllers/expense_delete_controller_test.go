package controllers

import (
	"net/http"
	"strconv"
	"testing"
)

func TestDeleteExpense(t *testing.T) {
	tests := []struct {
		name       string
		target     string
		seed       bool
		wantStatus int
		wantMsg    string
	}{
		{name: "valid delete", seed: true, wantStatus: http.StatusOK, wantMsg: "Expense deleted successfully"},
		{name: "missing expense", target: "/api/v1/expenses/99", wantStatus: http.StatusNotFound, wantMsg: "Expense not found"},
		{name: "invalid id", target: "/api/v1/expenses/abc", wantStatus: http.StatusBadRequest, wantMsg: "Invalid expense ID"},
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

			response := performRequest(http.MethodDelete, target, "", strconv.Itoa(user.ID))
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
