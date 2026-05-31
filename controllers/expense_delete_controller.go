package controllers

import (
	"errors"
	"net/http"

	"expense-tracker-api/models"

	beego "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

// ExpenseDeleteController handles expense deletion requests.
type ExpenseDeleteController struct {
	web.Controller
}

// Delete removes one expense owned by the authenticated user.
func (c *ExpenseDeleteController) Delete() {
	userID, ok := GetAuthenticatedUserID(c.Ctx)
	if !ok {
		beego.Warn("unauthorized expense delete attempt")
		return
	}

	expenseID, ok := parseExpenseID(&c.Controller)
	if !ok {
		return
	}

	err := models.DeleteExpense(expenseID, userID)
	if errors.Is(err, models.ErrExpenseNotFound) {
		writeExpenseError(&c.Controller, http.StatusNotFound, "Expense not found")
		return
	}
	if err != nil {
		beego.Error("failed to delete expense:", err)
		writeExpenseError(&c.Controller, http.StatusInternalServerError, "Could not delete expense")
		return
	}

	beego.Info("expense deleted:", expenseID)
	writeExpenseJSON(&c.Controller, http.StatusOK, true, "Expense deleted successfully", nil)
}
