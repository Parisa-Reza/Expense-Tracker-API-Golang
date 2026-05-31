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
func (c *ExpenseUpdateController) Update() {
	userID, ok := GetAuthenticatedUserID(c.Ctx)
	if !ok {
		beego.Warn("unauthorized expense update attempt")
		return
	}

	expenseID, ok := parseExpenseID(&c.Controller)
	if !ok {
		return
	}

	existingExpense, err := models.GetExpenseByID(expenseID, userID)
	if errors.Is(err, models.ErrExpenseNotFound) {
		writeExpenseError(&c.Controller, http.StatusNotFound, "Expense not found")
		return
	}
	if err != nil {
		beego.Error("failed to load expense for update:", err)
		writeExpenseError(&c.Controller, http.StatusInternalServerError, "Could not update expense")
		return
	}

	request, ok := parseExpenseRequest(&c.Controller)
	if !ok {
		return
	}

	existingExpense.Title = request.Title
	existingExpense.Amount = request.Amount
	existingExpense.Category = request.Category
	existingExpense.Note = request.Note
	existingExpense.ExpenseDate = request.ExpenseDate

	if err := models.UpdateExpense(existingExpense); errors.Is(err, models.ErrInvalidExpenseCategory) {
		writeExpenseError(&c.Controller, http.StatusBadRequest, "Invalid expense category")
		return
	} else if err != nil {
		beego.Error("failed to update expense:", err)
		writeExpenseError(&c.Controller, http.StatusInternalServerError, "Could not update expense")
		return
	}

	beego.Info("expense updated:", expenseID)
	writeExpenseJSON(&c.Controller, http.StatusOK, true, "Expense updated successfully", newExpenseResponse(*existingExpense))
}
