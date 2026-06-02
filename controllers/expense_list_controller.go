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
// @Title       List Expenses
// @Summary     Retrieves all expenses for authenticated user
// @Description Returns a list of all expenses for the authenticated user with optional filtering and sorting
// @Tags        Expenses
// @Produce     json
// @Param       X-User-ID header string true "User ID from login"
// @Param       category query string false "Filter by category"
// @Param       sort query string false "Sort by field (e.g., 'amount', 'date')"
// @Param       order query string false "Sort order: 'asc' or 'desc'"
// @Param       limit query int false "Maximum number of expenses to return"
// @Success     200 {object} map[string]interface{} "Expenses retrieved"
// @Failure     400 {object} map[string]interface{} "Invalid query parameters"
// @Failure     401 {object} map[string]interface{} "Unauthorized - User ID not provided"
// @Failure     500 {object} map[string]interface{} "Internal server error"
// @Security    UserIDHeader
// @router      /expenses [get]
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
