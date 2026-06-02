package controllers

import (
	"net/http"
	"strings"
	"time"

	"expense-tracker-api/models"

	"github.com/beego/beego/v2/server/web"
)

const (
	defaultExpenseSortOrder = "desc"
	expenseDateLayout       = "2006-01-02"
)

// parseExpenseListOptions extracts and validates query parameters for listing expenses.
// It ensures proper filtering (category, date range) and sorting rules.
// Returns a validated ExpenseListOptions struct and a boolean indicating success/failure.
func parseExpenseListOptions(controller *web.Controller) (models.ExpenseListOptions, bool) {
	options := models.ExpenseListOptions{
		Category:  strings.TrimSpace(controller.GetString("category")),
		DateFrom:  strings.TrimSpace(controller.GetString("date_from")),
		DateTo:    strings.TrimSpace(controller.GetString("date_to")),
		SortBy:    strings.TrimSpace(controller.GetString("sort_by")),
		SortOrder: strings.TrimSpace(controller.GetString("sort_order")),
	}

	// Apply default sorting order if not explicitly provided
	if options.SortOrder == "" {
		options.SortOrder = defaultExpenseSortOrder
	}

	// Validate date query parameters (format + optional values allowed)
	if !validateExpenseQueryDate(controller, "date_from", options.DateFrom) {
		return models.ExpenseListOptions{}, false
	}

	if !validateExpenseQueryDate(controller, "date_to", options.DateTo) {
		return models.ExpenseListOptions{}, false
	}

	// Ensure logical consistency of date range (from <= to)
	if !validateExpenseDateRange(controller, options.DateFrom, options.DateTo) {
		return models.ExpenseListOptions{}, false
	}

	// Restrict sorting field to supported columns only
	if options.SortBy != "" && options.SortBy != "amount" && options.SortBy != "expense_date" {
		writeExpenseError(controller, http.StatusBadRequest, "sort_by must be amount or expense_date")
		return models.ExpenseListOptions{}, false
	}

	// Validate sort order direction (prevents invalid query values)
	if options.SortOrder != "asc" && options.SortOrder != "desc" {
		writeExpenseError(controller, http.StatusBadRequest, "sort_order must be asc or desc")
		return models.ExpenseListOptions{}, false
	}

	return options, true
}

// parseExpenseSummaryRange extracts and validates optional date range used specifically for summary/aggregation endpoints.
func parseExpenseSummaryRange(controller *web.Controller) (string, string, bool) {
	dateFrom := strings.TrimSpace(controller.GetString("date_from"))
	dateTo := strings.TrimSpace(controller.GetString("date_to"))

	// Validate date format for start boundary (if provided)
	if !validateExpenseQueryDate(controller, "date_from", dateFrom) {
		return "", "", false
	}

	// Validate date format for end boundary (if provided)
	if !validateExpenseQueryDate(controller, "date_to", dateTo) {
		return "", "", false
	}

	// Ensure the range is logically valid
	if !validateExpenseDateRange(controller, dateFrom, dateTo) {
		return "", "", false
	}

	return dateFrom, dateTo, true
}

// validateExpenseQueryDate ensures a single date query parameter is either:
// 1. empty (optional parameter), or 2. in valid YYYY-MM-DD format

func validateExpenseQueryDate(controller *web.Controller, name string, value string) bool {

	// Empty values are allowed (parameter is optional)
	if value == "" {
		return true
	}

	// Strict validation of expected date format
	if _, err := time.Parse(expenseDateLayout, value); err != nil {
		writeExpenseError(controller, http.StatusBadRequest, name+" must use YYYY-MM-DD format")
		return false
	}

	return true
}

// validateExpenseDateRange ensures the provided date range is logically valid. It assumes both dates are already format-validated.
func validateExpenseDateRange(controller *web.Controller, dateFrom string, dateTo string) bool {

	// Only validate when both boundaries are provided . Lexicographical comparison works because format is YYYY-MM-DD

	if dateFrom != "" && dateTo != "" && dateFrom > dateTo {
		writeExpenseError(controller, http.StatusBadRequest, "date_from must be on or before date_to")
		return false
	}

	return true
}
