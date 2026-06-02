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

	// Validate user and get authenticated user ID.
	userID, ok := GetAuthenticatedUserID(c.Ctx)
	if !ok {
		beego.Warn("unauthorized expense list attempt")
		return
	}

	// Fetches only expenses belonging to user ID
	expenses, err := models.GetExpensesByUserID(userID)
	if err != nil {
		beego.Error("failed to list expenses:", err)

		// Return HTTP 500 if expenses cannot be loaded.
		writeExpenseError(&c.Controller, http.StatusInternalServerError, "Could not retrieve expenses")
		return
	}

	// Read filtering and sorting options from query parameters.
	options, ok := parseExpenseListOptions(&c.Controller)
	if !ok {
		return
	}

	// Apply filters and sorting.
	expenses = models.FilterAndSortExpenses(expenses, options)

	// Read and validate limit query parameter.
	limit, ok := parseLimit(&c.Controller)
	if !ok {
		return
	}

	// Keep only the first N expenses if a limit is provided.
	if limit > 0 && limit < len(expenses) {
		expenses = expenses[:limit]
	}

	// Prepare API response objects.
	data := make([]expenseResponse, 0, len(expenses))

	// Convert model expenses into response objects.
	for _, expense := range expenses {
		data = append(data, newExpenseResponse(expense))
	}

	beego.Info("expenses retrieved for user:", userID)

	// Return success response.
	writeExpenseJSON(&c.Controller, http.StatusOK, true, "Expenses retrieved", data)
}

// parseLimit reads and validates the "limit" query parameter.
func parseLimit(controller *web.Controller) (int, bool) {

	// Read limit parameter and remove surrounding spaces.
	limitValue := strings.TrimSpace(controller.GetString("limit"))

	// If limit is empty, return 0 (no limit) and true (valid).
	if limitValue == "" {
		return 0, true
	}

	// Convert string value to integer.
	limit, err := strconv.Atoi(limitValue)

	// Limit must be a positive integer.
	if err != nil || limit < 1 {
		writeExpenseError(controller, http.StatusBadRequest, "limit must be a positive integer")
		return 0, false
	}

	// Return validated limit.
	return limit, true
}
