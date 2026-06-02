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
// @Title       Delete Expense
// @Summary     Deletes an expense
// @Description Removes an expense from the authenticated user's records
// @Tags        Expenses
// @Produce     json
// @Param       X-User-ID header string true "User ID from login"
// @Param       id path int true "Expense ID"
// @Success     200 {object} map[string]interface{} "Expense deleted successfully"
// @Failure     400 {object} map[string]interface{} "Invalid expense ID"
// @Failure     401 {object} map[string]interface{} "Unauthorized - User ID not provided"
// @Failure     404 {object} map[string]interface{} "Expense not found"
// @Failure     500 {object} map[string]interface{} "Internal server error"
// @Security    UserIDHeader
// @Router      /expenses/{id} [delete]
func (c *ExpenseDeleteController) Delete() {

	// Get authenticated user's ID from request header.
	userID, ok := GetAuthenticatedUserID(c.Ctx)
	if !ok {
		beego.Warn("unauthorized expense delete attempt")
		return
	}

	// Get expense ID from URL parameter.
	expenseID, ok := parseExpenseID(&c.Controller)
	if !ok {
		return
	}

	// Attempt to delete the expense from CSV.
	err := models.DeleteExpense(expenseID, userID)

	// Handle case where expense does not exist.
	if errors.Is(err, models.ErrExpenseNotFound) {
		writeExpenseError(&c.Controller, http.StatusNotFound, "Expense not found")
		return
	}
	if err != nil {
		beego.Error("failed to delete expense:", err)

		// Handle unexpected internal errors.
		writeExpenseError(&c.Controller, http.StatusInternalServerError, "Could not delete expense")
		return
	}

	beego.Info("expense deleted:", expenseID)

	// Return success response.
	writeExpenseJSON(&c.Controller, http.StatusOK, true, "Expense deleted successfully", nil)
}
