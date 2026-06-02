package controllers

import (
	"net/http"
	"strconv"
	"testing"
)

func TestExpenseQueryParameters(t *testing.T) {
	tests := []struct {
		name       string
		target     string
		wantStatus int
		wantMsg    string
	}{
		{name: "valid sort", target: "/api/v1/expenses?sort_by=amount&sort_order=asc", wantStatus: http.StatusOK, wantMsg: "Expenses retrieved"},
		{name: "invalid sort by", target: "/api/v1/expenses?sort_by=title", wantStatus: http.StatusBadRequest, wantMsg: "sort_by must be amount or expense_date"},
		{name: "invalid sort order", target: "/api/v1/expenses?sort_order=sideways", wantStatus: http.StatusBadRequest, wantMsg: "sort_order must be asc or desc"},
		{name: "invalid date from", target: "/api/v1/expenses?date_from=2025/06/01", wantStatus: http.StatusBadRequest, wantMsg: "date_from must use YYYY-MM-DD format"},
		{name: "invalid date range", target: "/api/v1/expenses?date_from=2025-06-30&date_to=2025-06-01", wantStatus: http.StatusBadRequest, wantMsg: "date_from must be on or before date_to"},
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
