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
	userID, ok := GetAuthenticatedUserID(c.Ctx)
	if !ok {
		beego.Warn("unauthorized expense summary attempt")
		return
	}

	dateFrom, dateTo, ok := parseExpenseSummaryRange(&c.Controller)
	if !ok {
		return
	}

	summary, err := models.GetExpenseSummaryByUserID(userID, dateFrom, dateTo)
	if err != nil {
		beego.Error("failed to generate expense summary:", err)
		writeExpenseError(&c.Controller, http.StatusInternalServerError, "Could not generate summary")
		return
	}

	beego.Info("expense summary generated for user:", userID)
	writeExpenseJSON(&c.Controller, http.StatusOK, true, "Summary generated", summary)
}
