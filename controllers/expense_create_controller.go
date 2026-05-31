package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"expense-tracker-api/models"

	beego "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

// ExpenseCreateController handles expense creation requests.
type ExpenseCreateController struct {
	web.Controller
}

type expenseRequest struct {
	Title       string  `json:"title"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Note        string  `json:"note"`
	ExpenseDate string  `json:"expense_date"`
}

type expenseResponse struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Note        string  `json:"note"`
	ExpenseDate string  `json:"expense_date"`
}

// Create creates a new expense for the authenticated user.
func (c *ExpenseCreateController) Create() {
	userID, ok := GetAuthenticatedUserID(c.Ctx)
	if !ok {
		beego.Warn("unauthorized expense create attempt")
		return
	}

	request, ok := parseExpenseRequest(&c.Controller)
	if !ok {
		return
	}

	expense := &models.Expense{
		UserID:      userID,
		Title:       request.Title,
		Amount:      request.Amount,
		Category:    request.Category,
		Note:        request.Note,
		ExpenseDate: request.ExpenseDate,
	}

	if err := models.CreateExpense(expense); errors.Is(err, models.ErrInvalidExpenseCategory) {
		writeExpenseError(&c.Controller, http.StatusBadRequest, "Invalid expense category")
		return
	} else if err != nil {
		beego.Error("failed to create expense:", err)
		writeExpenseError(&c.Controller, http.StatusInternalServerError, "Could not create expense")
		return
	}

	beego.Info("expense created:", expense.ID)
	writeExpenseJSON(&c.Controller, http.StatusCreated, true, "Expense created successfully", newExpenseResponse(*expense))
}

func parseExpenseRequest(controller *web.Controller) (expenseRequest, bool) {
	var request expenseRequest
	if err := json.Unmarshal(controller.Ctx.Input.RequestBody, &request); err != nil {
		beego.Warn("invalid expense request body:", err)
		writeExpenseError(controller, http.StatusBadRequest, "Invalid request body")
		return expenseRequest{}, false
	}

	request.Title = strings.TrimSpace(request.Title)
	request.Category = strings.TrimSpace(request.Category)
	request.Note = strings.TrimSpace(request.Note)
	request.ExpenseDate = strings.TrimSpace(request.ExpenseDate)

	if request.Title == "" || request.Category == "" || request.ExpenseDate == "" {
		writeExpenseError(controller, http.StatusBadRequest, "Title, category, and expense_date are required")
		return expenseRequest{}, false
	}

	if request.Amount <= 0 {
		writeExpenseError(controller, http.StatusBadRequest, "Amount must be greater than zero")
		return expenseRequest{}, false
	}

	if !models.IsAllowedCategory(request.Category) {
		writeExpenseError(controller, http.StatusBadRequest, "Invalid expense category")
		return expenseRequest{}, false
	}

	if _, err := time.Parse("2006-01-02", request.ExpenseDate); err != nil {
		writeExpenseError(controller, http.StatusBadRequest, "expense_date must use YYYY-MM-DD format")
		return expenseRequest{}, false
	}

	return request, true
}

func newExpenseResponse(expense models.Expense) expenseResponse {
	return expenseResponse{
		ID:          expense.ID,
		Title:       expense.Title,
		Amount:      expense.Amount,
		Category:    expense.Category,
		Note:        expense.Note,
		ExpenseDate: expense.ExpenseDate,
	}
}

func writeExpenseJSON(controller *web.Controller, statusCode int, success bool, message string, data interface{}) {
	controller.Ctx.Output.SetStatus(statusCode)
	controller.Data["json"] = map[string]interface{}{
		"success": success,
		"message": message,
		"data":    data,
	}
	controller.ServeJSON()
}

func writeExpenseError(controller *web.Controller, statusCode int, message string) {
	controller.Ctx.Output.SetStatus(statusCode)
	controller.Data["json"] = map[string]interface{}{
		"success": false,
		"message": message,
		"data":    nil,
	}
	controller.ServeJSON()
}
