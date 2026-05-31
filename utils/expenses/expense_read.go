package expenseutils

import (
	"encoding/csv"
	"os"
)

// ReadExpensesCSV reads all expense data rows from the expenses CSV file.
func ReadExpensesCSV() ([][]string, error) {
	if err := EnsureExpensesCSV(); err != nil {
		return nil, err
	}

	file, err := os.Open(GetExpensesCSVPath())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) <= 1 {
		return [][]string{}, nil
	}

	return records[1:], nil
}
