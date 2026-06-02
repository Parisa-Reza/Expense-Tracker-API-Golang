package expenseutils

import (
	"encoding/csv"
	"os"
	"path/filepath"

	beego "github.com/beego/beego/v2/server/web"
)

// CSV header row
var expensesCSVHeader = []string{"id", "user_id", "title", "amount", "category", "note", "expense_date", "created_at"}

// GetExpensesCSVPath returns the configured expenses CSV path from app.conf.
func GetExpensesCSVPath() string {
	return beego.AppConfig.DefaultString("expenses_csv_path", "")
}

// EnsureExpensesCSV creates the expenses CSV file with a header when it does not exist.
func EnsureExpensesCSV() error {
	expensesCSVPath := GetExpensesCSVPath()

	// Create folder if missing

	// here 0755 means: read/write permissions for owner, and read permissions for group and others

	if err := os.MkdirAll(filepath.Dir(expensesCSVPath), 0755); err != nil {
		return err
	}

	// Check if file already exists
	if _, err := os.Stat(expensesCSVPath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	// Create the file
	file, err := os.Create(expensesCSVPath)
	if err != nil {
		return err
	}

	// Close the file when function ends
	defer file.Close()

	// : Write header row
	writer := csv.NewWriter(file)
	if err := writer.Write(expensesCSVHeader); err != nil {
		return err
	}

	// Flush() forces buffered (temporary) CSV data to be permanently written into the file.
	writer.Flush()

	// If anything went wrong while writing, return the error. Otherwise returns nil
	return writer.Error()
}

// // AppendExpenseCSV appends one expense record to the expenses CSV file.
// func AppendExpenseCSV(record []string) error {

// 	// Make sure the CSV file exists (create it if not)
// 	if err := EnsureExpensesCSV(); err != nil {
// 		return err
// 	}

// 	// Open the file in append mode so we can add new data at the end . 0644 means: read/write permissions for owner, and read permissions for group and others
// 	file, err := os.OpenFile(GetExpensesCSVPath(), os.O_APPEND|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return err
// 	}

// 	// Close file when function finishes
// 	defer file.Close()

// 	// Create a CSV writer to write data into the file
// 	writer := csv.NewWriter(file)

// 	// Write one row (user record) into CSV
// 	if err := writer.Write(record); err != nil {
// 		return err
// 	}

// 	// Flush ensures data is written from memory buffer to file
// 	writer.Flush()

// 	// Return any error that happened during writing or flushing
// 	return writer.Error()
// }
