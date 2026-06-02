package controllers

import (
	"net/http"

	"expense-tracker-api/models"

	beego "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

// ExpenseSummaryController handles expense summary requests.
type ExpenseSummaryController struct {
	web.Controller
}

// Summary returns aggregate expense totals for the authenticated user.
// @Title       Get Expense Summary
// @Summary     Retrieves expense summary and totals
// @Description Returns aggregated expense totals by category and overall for the authenticated user within optional date range
// @Tags        Expenses
// @Produce     json
// @Param       X-User-ID header string true "User ID from login"
// @Param       date_from query string false "Start date (YYYY-MM-DD format)"
// @Param       date_to query string false "End date (YYYY-MM-DD format)"
// @Success     200 {object} map[string]interface{} "Summary generated"
// @Failure     400 {object} map[string]interface{} "Invalid date format or query parameters"
// @Failure     401 {object} map[string]interface{} "Unauthorized - User ID not provided"
// @Failure     500 {object} map[string]interface{} "Internal server error"
// @Security    UserIDHeader
// @Router      /expenses/summary [get]
func (c *ExpenseSummaryController) Summary() {

	// Extract authenticated user ID from request context/header.
	userID, ok := GetAuthenticatedUserID(c.Ctx)
	if !ok {
		beego.Warn("unauthorized expense summary attempt")
		return
	}

	// Parse date range (start and end) from request parameters.
	dateFrom, dateTo, ok := parseExpenseSummaryRange(&c.Controller)
	if !ok {
		return
	}

	// Fetch aggregated expense summary for the user within the given date range.
	summary, err := models.GetExpenseSummaryByUserID(userID, dateFrom, dateTo)

	if err != nil {
		beego.Error("failed to generate expense summary:", err)
		writeExpenseError(&c.Controller, http.StatusInternalServerError, "Could not generate summary")
		return
	}

	beego.Info("expense summary generated for user:", userID)

	// Return aggregated summary response.
	writeExpenseJSON(&c.Controller, http.StatusOK, true, "Summary generated", summary)
}
