package controllers

import (
	"net/http"
	"strconv"
	"testing"
)

func TestExpenseSummary(t *testing.T) {
	tests := []struct {
		name       string
		target     string
		wantStatus int
		wantMsg    string
	}{
		{name: "valid summary", target: "/api/v1/expenses/summary", wantStatus: http.StatusOK, wantMsg: "Summary generated"},
		{name: "valid date range", target: "/api/v1/expenses/summary?date_from=2025-06-01&date_to=2025-06-30", wantStatus: http.StatusOK, wantMsg: "Summary generated"},
		{name: "invalid date", target: "/api/v1/expenses/summary?date_from=06-01-2025", wantStatus: http.StatusBadRequest, wantMsg: "date_from must use YYYY-MM-DD format"},
		{name: "invalid range", target: "/api/v1/expenses/summary?date_from=2025-06-30&date_to=2025-06-01", wantStatus: http.StatusBadRequest, wantMsg: "date_from must be on or before date_to"},
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
