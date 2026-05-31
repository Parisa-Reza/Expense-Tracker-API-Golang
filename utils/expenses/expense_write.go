package expenseutils

import (
	"encoding/csv"
	"os"
)

// WriteExpensesCSV rewrites the entire expenses CSV file with the provided data rows.
func WriteExpensesCSV(records [][]string) error {
	if err := EnsureExpensesCSV(); err != nil {
		return err
	}

	file, err := os.Create(GetExpensesCSVPath())
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.Write(expensesCSVHeader); err != nil {
		return err
	}

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	writer.Flush()
	return writer.Error()
}
