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
// @Title       Get Expense
// @Summary     Retrieves a single expense by ID
// @Description Returns details of a specific expense owned by the authenticated user
// @Tags        Expenses
// @Produce     json
// @Param       X-User-ID header string true "User ID from login"
// @Param       id path int true "Expense ID"
// @Success     200 {object} map[string]interface{} "Expense retrieved"
// @Failure     400 {object} map[string]interface{} "Invalid expense ID"
// @Failure     401 {object} map[string]interface{} "Unauthorized - User ID not provided"
// @Failure     404 {object} map[string]interface{} "Expense not found"
// @Failure     500 {object} map[string]interface{} "Internal server error"
// @Security    UserIDHeader
// @router      /expenses/{id} [get]
func (c *ExpenseGetController) Get() {

	// GetAuthenticatedUserID to checks if the request is authenticated and to get the user ID.
	userID, ok := GetAuthenticatedUserID(c.Ctx)
	if !ok {
		beego.Warn("unauthorized expense get attempt")
		return
	}

	expenseID, ok := parseExpenseID(&c.Controller)
	if !ok {
		return
	}

	// Fetch the expense, scoping the lookup to the authenticated user so that one user cannot read another user's records.
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

	// here writeExpensionJSON to send a JSON response back to the client with the expense data, using the newExpenseResponse helper function to format the expense data appropriately for the response.
	writeExpenseJSON(&c.Controller, http.StatusOK, true, "Expense retrieved", newExpenseResponse(*expense))
}

// parseExpenseID extracts the ":id" route parameter from the request URL, converts it to an integer, and validates that it is a positive value.
func parseExpenseID(controller *web.Controller) (int, bool) {
	expenseID, err := strconv.Atoi(controller.Ctx.Input.Param(":id"))
	if err != nil || expenseID < 1 {
		writeExpenseError(controller, http.StatusBadRequest, "Invalid expense ID")
		return 0, false
	}

	return expenseID, true
}
