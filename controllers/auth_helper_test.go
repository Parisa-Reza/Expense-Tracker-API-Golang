package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"

	"expense-tracker-api/models"

	"github.com/beego/beego/v2/server/web"
)

type apiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var testRoutesOnce sync.Once

func registerTestRoutes() {
	testRoutesOnce.Do(func() {
		web.BConfig.CopyRequestBody = true
		web.Router("/api/v1/health", &HealthController{}, "get:Check")
		web.Router("/api/v1/auth/register", &RegisterController{}, "post:Register")
		web.Router("/api/v1/auth/login", &LoginController{}, "post:Login")
		web.Router("/api/v1/expenses", &ExpenseCreateController{}, "post:Create")
		web.Router("/api/v1/expenses", &ExpenseListController{}, "get:List")
		web.Router("/api/v1/expenses/summary", &ExpenseSummaryController{}, "get:Summary")
		web.Router("/api/v1/expenses/:id", &ExpenseGetController{}, "get:Get")
		web.Router("/api/v1/expenses/:id", &ExpenseUpdateController{}, "put:Update")
		web.Router("/api/v1/expenses/:id", &ExpenseDeleteController{}, "delete:Delete")
	})
}

func setupControllerTest(t *testing.T) {
	t.Helper()
	registerTestRoutes()
	dir := t.TempDir()
	if err := web.AppConfig.Set("users_csv_path", filepath.Join(dir, "users.csv")); err != nil {
		t.Fatalf("set users csv path: %v", err)
	}
	if err := web.AppConfig.Set("expenses_csv_path", filepath.Join(dir, "expenses.csv")); err != nil {
		t.Fatalf("set expenses csv path: %v", err)
	}
}

func performRequest(method string, target string, body string, userID string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	if userID != "" {
		request.Header.Set("X-User-ID", userID)
	}
	response := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(response, request)
	return response
}

func decodeAPIResponse(t *testing.T, response *httptest.ResponseRecorder) apiResponse {
	t.Helper()
	var got apiResponse
	if err := json.Unmarshal(response.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response %q: %v", response.Body.String(), err)
	}
	return got
}

func createTestUser(t *testing.T) models.User {
	t.Helper()
	user := models.User{Name: "Test User", Email: "test@example.com", Password: "secret1"}
	if err := models.CreateUser(&user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	return user
}

func createTestExpense(t *testing.T, userID int) models.Expense {
	t.Helper()
	expense := models.Expense{
		UserID:      userID,
		Title:       "Lunch",
		Amount:      350.50,
		Category:    "Food",
		Note:        "team meal",
		ExpenseDate: "2025-06-10",
	}
	if err := models.CreateExpense(&expense); err != nil {
		t.Fatalf("create expense: %v", err)
	}
	return expense
}

func TestRequireUserIDThroughExpenseEndpoint(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		wantStatus int
		wantMsg    string
	}{
		{name: "missing user id", userID: "", wantStatus: http.StatusUnauthorized, wantMsg: "Unauthorized"},
		{name: "invalid user id", userID: "abc", wantStatus: http.StatusUnauthorized, wantMsg: "Unauthorized"},
		{name: "unknown user id", userID: "99", wantStatus: http.StatusUnauthorized, wantMsg: "Unauthorized"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupControllerTest(t)
			response := performRequest(http.MethodGet, "/api/v1/expenses", "", tt.userID)
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
