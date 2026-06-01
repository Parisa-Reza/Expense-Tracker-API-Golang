package expenseutils

import (
	"encoding/csv"
	"os"
)

// ReadExpensesCSV reads all expense data rows from the expenses CSV file.
func ReadExpensesCSV() ([][]string, error) {

	// Ensure the CSV file exists
	if err := EnsureExpensesCSV(); err != nil {
		return nil, err
	}


    // Open the CSV file for reading
	file, err := os.Open(GetExpensesCSVPath())
	if err != nil {
		return nil, err
	}

	// Ensure file is closed after function ends
	defer file.Close()

	// Read all rows from CSV file into memory
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	// If file has only header or is empty, return empty result
	if len(records) <= 1 {
		return [][]string{}, nil
	}

	// Skip the first row (header) and return only user data
	return records[1:], nil
}
