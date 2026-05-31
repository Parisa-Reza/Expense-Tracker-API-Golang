package expenseutils

import (
	"encoding/csv"
	"os"
	"path/filepath"

	beego "github.com/beego/beego/v2/server/web"
)

// ExpensesCSVPath is the fallback path for expense CSV storage.
const ExpensesCSVPath = "data/expenses.csv"

var expensesCSVHeader = []string{"id", "user_id", "title", "amount", "category", "note", "expense_date", "created_at"}

// GetExpensesCSVPath returns the configured expenses CSV path.
func GetExpensesCSVPath() string {
	return beego.AppConfig.DefaultString("expenses_csv_path", ExpensesCSVPath)
}

// EnsureExpensesCSV creates the expenses CSV file with a header when it does not exist.
func EnsureExpensesCSV() error {
	expensesCSVPath := GetExpensesCSVPath()

	if err := os.MkdirAll(filepath.Dir(expensesCSVPath), 0755); err != nil {
		return err
	}

	if _, err := os.Stat(expensesCSVPath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	file, err := os.Create(expensesCSVPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.Write(expensesCSVHeader); err != nil {
		return err
	}

	writer.Flush()
	return writer.Error()
}

// AppendExpenseCSV appends one expense record to the expenses CSV file.
func AppendExpenseCSV(record []string) error {
	if err := EnsureExpensesCSV(); err != nil {
		return err
	}

	file, err := os.OpenFile(GetExpensesCSVPath(), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.Write(record); err != nil {
		return err
	}

	writer.Flush()
	return writer.Error()
}
