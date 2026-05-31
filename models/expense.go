package models

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	expenseutils "expense-tracker-api/utils/expenses"
)

// AllowedCategories contains every valid expense category.
var AllowedCategories = []string{
	"Food", "Transport", "Housing", "Entertainment", "Shopping", "Healthcare", "Education", "Utilities", "Other",
}

// Expense represents one expense record stored in CSV.
type Expense struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Note        string    `json:"note"`
	ExpenseDate string    `json:"expense_date"`
	CreatedAt   time.Time `json:"created_at"`
}

var (
	// ErrExpenseNotFound is returned when an expense cannot be found for the requested user.
	ErrExpenseNotFound = errors.New("expense not found")
	// ErrInvalidExpenseCategory is returned when an expense category is not allowed.
	ErrInvalidExpenseCategory = errors.New("invalid expense category")
)

// ExpenseListOptions contains optional filters and sorting for expense lists.
type ExpenseListOptions struct {
	Category  string
	DateFrom  string
	DateTo    string
	SortBy    string
	SortOrder string
}

// ExpenseSummary contains aggregate expense data for a date range.
type ExpenseSummary struct {
	DateFrom       string             `json:"date_from"`
	DateTo         string             `json:"date_to"`
	TotalAmount    float64            `json:"total_amount"`
	TotalCount     int                `json:"total_count"`
	CategoryTotals map[string]float64 `json:"category_totals"`
}

// GetAllExpenses returns every valid expense row from the CSV storage.
func GetAllExpenses() ([]Expense, error) {
	records, err := expenseutils.ReadExpensesCSV()
	if err != nil {
		return nil, err
	}

	expenses := make([]Expense, 0, len(records))
	for _, record := range records {
		expense, err := expenseFromRecord(record)
		if err != nil {
			continue
		}
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

// GetExpensesByUserID returns all expenses owned by a user.
func GetExpensesByUserID(userID int) ([]Expense, error) {
	expenses, err := GetAllExpenses()
	if err != nil {
		return nil, err
	}

	userExpenses := make([]Expense, 0)
	for _, expense := range expenses {
		if expense.UserID == userID {
			userExpenses = append(userExpenses, expense)
		}
	}

	return userExpenses, nil
}

// GetExpensesByUserIDWithOptions returns a user's expenses after applying filters and sorting.
func GetExpensesByUserIDWithOptions(userID int, options ExpenseListOptions) ([]Expense, error) {
	expenses, err := GetExpensesByUserID(userID)
	if err != nil {
		return nil, err
	}

	return FilterAndSortExpenses(expenses, options), nil
}

// FilterAndSortExpenses applies category, date range, and sort options to expenses.
func FilterAndSortExpenses(expenses []Expense, options ExpenseListOptions) []Expense {
	filteredExpenses := make([]Expense, 0, len(expenses))
	for _, expense := range expenses {
		if options.Category != "" && expense.Category != options.Category {
			continue
		}
		if options.DateFrom != "" && expense.ExpenseDate < options.DateFrom {
			continue
		}
		if options.DateTo != "" && expense.ExpenseDate > options.DateTo {
			continue
		}

		filteredExpenses = append(filteredExpenses, expense)
	}

	sortExpenses(filteredExpenses, options.SortBy, options.SortOrder)
	return filteredExpenses
}

// GetExpenseSummaryByUserID returns aggregate totals for a user's expenses in a date range.
func GetExpenseSummaryByUserID(userID int, dateFrom string, dateTo string) (ExpenseSummary, error) {
	expenses, err := GetExpensesByUserID(userID)
	if err != nil {
		return ExpenseSummary{}, err
	}

	return SummarizeExpenses(expenses, dateFrom, dateTo), nil
}

// SummarizeExpenses totals expenses that fall within the optional date range.
func SummarizeExpenses(expenses []Expense, dateFrom string, dateTo string) ExpenseSummary {
	summary := ExpenseSummary{
		DateFrom:       dateFrom,
		DateTo:         dateTo,
		CategoryTotals: make(map[string]float64),
	}

	for _, expense := range expenses {
		if dateFrom != "" && expense.ExpenseDate < dateFrom {
			continue
		}
		if dateTo != "" && expense.ExpenseDate > dateTo {
			continue
		}

		summary.TotalAmount += expense.Amount
		summary.TotalCount++
		summary.CategoryTotals[expense.Category] += expense.Amount
	}

	return summary
}

// GetExpenseByID returns one expense when it belongs to the requested user.
func GetExpenseByID(id int, userID int) (*Expense, error) {
	expenses, err := GetExpensesByUserID(userID)
	if err != nil {
		return nil, err
	}

	for i := range expenses {
		if expenses[i].ID == id {
			return &expenses[i], nil
		}
	}

	return nil, ErrExpenseNotFound
}

// CreateExpense stores a new expense in the CSV storage.
func CreateExpense(expense *Expense) error {
	if !IsAllowedCategory(expense.Category) {
		return ErrInvalidExpenseCategory
	}

	expense.ID = GetNextExpenseID()
	expense.CreatedAt = time.Now().UTC()

	return expenseutils.AppendExpenseCSV(expenseToRecord(*expense))
}

// UpdateExpense rewrites the expense CSV with the changed expense row.
func UpdateExpense(expense *Expense) error {
	if !IsAllowedCategory(expense.Category) {
		return ErrInvalidExpenseCategory
	}

	expenses, err := GetAllExpenses()
	if err != nil {
		return err
	}

	found := false
	records := make([][]string, 0, len(expenses))
	for _, currentExpense := range expenses {
		if currentExpense.ID == expense.ID && currentExpense.UserID == expense.UserID {
			records = append(records, expenseToRecord(*expense))
			found = true
			continue
		}

		records = append(records, expenseToRecord(currentExpense))
	}

	if !found {
		return ErrExpenseNotFound
	}

	return expenseutils.WriteExpensesCSV(records)
}

// DeleteExpense rewrites the expense CSV without the deleted expense row.
func DeleteExpense(id int, userID int) error {
	expenses, err := GetAllExpenses()
	if err != nil {
		return err
	}

	found := false
	records := make([][]string, 0, len(expenses))
	for _, expense := range expenses {
		if expense.ID == id && expense.UserID == userID {
			found = true
			continue
		}

		records = append(records, expenseToRecord(expense))
	}

	if !found {
		return ErrExpenseNotFound
	}

	return expenseutils.WriteExpensesCSV(records)
}

// GetNextExpenseID returns the next available expense ID.
func GetNextExpenseID() int {
	expenses, err := GetAllExpenses()
	if err != nil || len(expenses) == 0 {
		return 1
	}

	maxID := 0
	for _, expense := range expenses {
		if expense.ID > maxID {
			maxID = expense.ID
		}
	}

	return maxID + 1
}

// IsAllowedCategory reports whether a category can be used for an expense.
func IsAllowedCategory(category string) bool {
	for _, allowedCategory := range AllowedCategories {
		if category == allowedCategory {
			return true
		}
	}

	return false
}

func sortExpenses(expenses []Expense, sortBy string, sortOrder string) {
	if sortBy == "" {
		return
	}

	descending := sortOrder == "desc"
	sort.Slice(expenses, func(i int, j int) bool {
		var comparison int
		switch sortBy {
		case "amount":
			if expenses[i].Amount < expenses[j].Amount {
				comparison = -1
			} else if expenses[i].Amount > expenses[j].Amount {
				comparison = 1
			}
		case "expense_date":
			if expenses[i].ExpenseDate < expenses[j].ExpenseDate {
				comparison = -1
			} else if expenses[i].ExpenseDate > expenses[j].ExpenseDate {
				comparison = 1
			}
		default:
			return false
		}

		if descending {
			return comparison > 0
		}

		return comparison < 0
	})
}

func expenseFromRecord(record []string) (Expense, error) {
	if len(record) < 8 {
		return Expense{}, errors.New("invalid expense record")
	}

	id, err := strconv.Atoi(record[0])
	if err != nil {
		return Expense{}, err
	}

	userID, err := strconv.Atoi(record[1])
	if err != nil {
		return Expense{}, err
	}

	amount, err := strconv.ParseFloat(record[3], 64)
	if err != nil {
		return Expense{}, err
	}

	createdAt, err := time.Parse(time.RFC3339, record[7])
	if err != nil {
		return Expense{}, err
	}

	return Expense{
		ID:          id,
		UserID:      userID,
		Title:       record[2],
		Amount:      amount,
		Category:    record[4],
		Note:        record[5],
		ExpenseDate: record[6],
		CreatedAt:   createdAt,
	}, nil
}

func expenseToRecord(expense Expense) []string {
	return []string{
		strconv.Itoa(expense.ID),
		strconv.Itoa(expense.UserID),
		strings.TrimSpace(expense.Title),
		strconv.FormatFloat(expense.Amount, 'f', 2, 64),
		expense.Category,
		strings.TrimSpace(expense.Note),
		expense.ExpenseDate,
		expense.CreatedAt.Format(time.RFC3339),
	}
}
