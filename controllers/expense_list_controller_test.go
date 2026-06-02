package controllers

import (
	"net/http"
	"strconv"
	"testing"
)

func TestListExpenses(t *testing.T) {
	tests := []struct {
		name       string
		target     string
		wantStatus int
		wantMsg    string
	}{
		{name: "all expenses", target: "/api/v1/expenses", wantStatus: http.StatusOK, wantMsg: "Expenses retrieved"},
		{name: "with limit", target: "/api/v1/expenses?limit=1", wantStatus: http.StatusOK, wantMsg: "Expenses retrieved"},
		{name: "invalid limit", target: "/api/v1/expenses?limit=bad", wantStatus: http.StatusBadRequest, wantMsg: "limit must be a positive integer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupControllerTest(t)
			user := createTestUser(t)
			createTestExpense(t, user.ID)

			response := performRequest(http.MethodGet, tt.target, "", strconv.Itoa(user.ID))
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
