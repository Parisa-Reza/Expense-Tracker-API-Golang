package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"expense-tracker-api/models"

	beego "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

// ExpenseGetController handles single expense retrieval requests.
type ExpenseGetController struct {
	web.Controller
}

// Get returns one expense owned by the authenticated user.
func (c *ExpenseGetController) Get() {
	userID, ok := GetAuthenticatedUserID(c.Ctx)
	if !ok {
		beego.Warn("unauthorized expense get attempt")
		return
	}

	expenseID, ok := parseExpenseID(&c.Controller)
	if !ok {
		return
	}

	expense, err := models.GetExpenseByID(expenseID, userID)
	if errors.Is(err, models.ErrExpenseNotFound) {
		writeExpenseError(&c.Controller, http.StatusNotFound, "Expense not found")
		return
	}
	if err != nil {
		beego.Error("failed to get expense:", err)
		writeExpenseError(&c.Controller, http.StatusInternalServerError, "Could not retrieve expense")
		return
	}

	beego.Info("expense retrieved:", expenseID)
	writeExpenseJSON(&c.Controller, http.StatusOK, true, "Expense retrieved", newExpenseResponse(*expense))
}

func parseExpenseID(controller *web.Controller) (int, bool) {
	expenseID, err := strconv.Atoi(controller.Ctx.Input.Param(":id"))
	if err != nil || expenseID < 1 {
		writeExpenseError(controller, http.StatusBadRequest, "Invalid expense ID")
		return 0, false
	}

	return expenseID, true
}
