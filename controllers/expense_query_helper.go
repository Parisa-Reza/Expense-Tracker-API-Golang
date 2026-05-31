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

func parseExpenseListOptions(controller *web.Controller) (models.ExpenseListOptions, bool) {
	options := models.ExpenseListOptions{
		Category:  strings.TrimSpace(controller.GetString("category")),
		DateFrom:  strings.TrimSpace(controller.GetString("date_from")),
		DateTo:    strings.TrimSpace(controller.GetString("date_to")),
		SortBy:    strings.TrimSpace(controller.GetString("sort_by")),
		SortOrder: strings.TrimSpace(controller.GetString("sort_order")),
	}

	if options.SortOrder == "" {
		options.SortOrder = defaultExpenseSortOrder
	}

	if !validateExpenseQueryDate(controller, "date_from", options.DateFrom) {
		return models.ExpenseListOptions{}, false
	}
	if !validateExpenseQueryDate(controller, "date_to", options.DateTo) {
		return models.ExpenseListOptions{}, false
	}
	if !validateExpenseDateRange(controller, options.DateFrom, options.DateTo) {
		return models.ExpenseListOptions{}, false
	}
	if options.SortBy != "" && options.SortBy != "amount" && options.SortBy != "expense_date" {
		writeExpenseError(controller, http.StatusBadRequest, "sort_by must be amount or expense_date")
		return models.ExpenseListOptions{}, false
	}
	if options.SortOrder != "asc" && options.SortOrder != "desc" {
		writeExpenseError(controller, http.StatusBadRequest, "sort_order must be asc or desc")
		return models.ExpenseListOptions{}, false
	}

	return options, true
}

func parseExpenseSummaryRange(controller *web.Controller) (string, string, bool) {
	dateFrom := strings.TrimSpace(controller.GetString("date_from"))
	dateTo := strings.TrimSpace(controller.GetString("date_to"))

	if !validateExpenseQueryDate(controller, "date_from", dateFrom) {
		return "", "", false
	}
	if !validateExpenseQueryDate(controller, "date_to", dateTo) {
		return "", "", false
	}
	if !validateExpenseDateRange(controller, dateFrom, dateTo) {
		return "", "", false
	}

	return dateFrom, dateTo, true
}

func validateExpenseQueryDate(controller *web.Controller, name string, value string) bool {
	if value == "" {
		return true
	}

	if _, err := time.Parse(expenseDateLayout, value); err != nil {
		writeExpenseError(controller, http.StatusBadRequest, name+" must use YYYY-MM-DD format")
		return false
	}

	return true
}

func validateExpenseDateRange(controller *web.Controller, dateFrom string, dateTo string) bool {
	if dateFrom != "" && dateTo != "" && dateFrom > dateTo {
		writeExpenseError(controller, http.StatusBadRequest, "date_from must be on or before date_to")
		return false
	}

	return true
}
