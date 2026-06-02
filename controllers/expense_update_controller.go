package controllers

import (
	"errors"
	"net/http"

	"expense-tracker-api/models"

	beego "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

// ExpenseUpdateController handles expense update requests.
type ExpenseUpdateController struct {
	web.Controller
}

// Update updates one expense owned by the authenticated user.
// @Title       Update Expense
// @Summary     Updates an existing expense
// @Description Updates an expense with new title, amount, category, note, or date
// @Tags        Expenses
// @Accept      json
// @Produce     json
// @Param       X-User-ID header string true "User ID from login"
// @Param       id path int true "Expense ID"
// @Param       expense body expenseRequest true "Updated expense data"
// @Success     200 {object} map[string]interface{} "Expense updated successfully"
// @Failure     400 {object} map[string]interface{} "Invalid request body or expense ID"
// @Failure     401 {object} map[string]interface{} "Unauthorized - User ID not provided"
// @Failure     404 {object} map[string]interface{} "Expense not found"
// @Failure     500 {object} map[string]interface{} "Internal server error"
// @Security    UserIDHeader
// @router      /expenses/{id} [put]
func (c *ExpenseUpdateController) Update() {

	// Get authenticated user's ID from request header.
	userID, ok := GetAuthenticatedUserID(c.Ctx)
	if !ok {
		beego.Warn("unauthorized expense update attempt")
		return
	}

	// Get expense ID from URL parameter.
	expenseID, ok := parseExpenseID(&c.Controller)
	if !ok {
		return
	}

	// Load the expense that belongs to this user.
	existingExpense, err := models.GetExpenseByID(expenseID, userID)

	// Expense does not exist.
	if errors.Is(err, models.ErrExpenseNotFound) {
		writeExpenseError(&c.Controller, http.StatusNotFound, "Expense not found")
		return
	}

	// file loading error.
	if err != nil {
		beego.Error("failed to load expense for update:", err)
		writeExpenseError(&c.Controller, http.StatusInternalServerError, "Could not update expense")
		return
	}

	// Parse request body JSON.
	request, ok := parseExpenseRequest(&c.Controller)
	if !ok {
		return
	}

	existingExpense.Title = request.Title
	existingExpense.Amount = request.Amount
	existingExpense.Category = request.Category
	existingExpense.Note = request.Note
	existingExpense.ExpenseDate = request.ExpenseDate

	// Save updated expense.
	if err := models.UpdateExpense(existingExpense); errors.Is(err, models.ErrInvalidExpenseCategory) {
		writeExpenseError(&c.Controller, http.StatusBadRequest, "Invalid expense category")
		return
	} else if err != nil {

		// Unexpected save error.
		beego.Error("failed to update expense:", err)
		writeExpenseError(&c.Controller, http.StatusInternalServerError, "Could not update expense")
		return
	}

	beego.Info("expense updated:", expenseID)

	// Return success response.
	writeExpenseJSON(&c.Controller, http.StatusOK, true, "Expense updated successfully", newExpenseResponse(*existingExpense))
}
