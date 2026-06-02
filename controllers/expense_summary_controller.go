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
