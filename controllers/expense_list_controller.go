package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"expense-tracker-api/models"

	beego "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

// ExpenseListController handles expense list requests.
type ExpenseListController struct {
	web.Controller
}

// List returns expenses owned by the authenticated user.
func (c *ExpenseListController) List() {
	userID, ok := GetAuthenticatedUserID(c.Ctx)
	if !ok {
		beego.Warn("unauthorized expense list attempt")
		return
	}

	expenses, err := models.GetExpensesByUserID(userID)
	if err != nil {
		beego.Error("failed to list expenses:", err)
		writeExpenseError(&c.Controller, http.StatusInternalServerError, "Could not retrieve expenses")
		return
	}

	limit, ok := parseLimit(&c.Controller)
	if !ok {
		return
	}

	if limit > 0 && limit < len(expenses) {
		expenses = expenses[:limit]
	}

	data := make([]expenseResponse, 0, len(expenses))
	for _, expense := range expenses {
		data = append(data, newExpenseResponse(expense))
	}

	beego.Info("expenses retrieved for user:", userID)
	writeExpenseJSON(&c.Controller, http.StatusOK, true, "Expenses retrieved", data)
}

func parseLimit(controller *web.Controller) (int, bool) {
	limitValue := strings.TrimSpace(controller.GetString("limit"))
	if limitValue == "" {
		return 0, true
	}

	limit, err := strconv.Atoi(limitValue)
	if err != nil || limit < 1 {
		writeExpenseError(controller, http.StatusBadRequest, "limit must be a positive integer")
		return 0, false
	}

	return limit, true
}
