package expenseutils

import (
	"encoding/csv"
	"os"
)
// WriteExpensesCSV rewrites the entire expenses CSV file with the provided data rows.
func WriteExpensesCSV(records [][]string) error {

	// Make sure the CSV file exists
	if err := EnsureExpensesCSV(); err != nil {
		return err
	}

	// Create the file (this will overwrite existing file)
	file, err := os.Create(GetExpensesCSVPath())
	if err != nil {
		return err
	}

	// Close file when function finishes
	defer file.Close()

	// Create a CSV writer to write data into the file
	writer := csv.NewWriter(file)

	// Write header row first.
	if err := writer.Write(expensesCSVHeader); err != nil {
		return err
	}

	// write all expense records. If any error happens during writing, return it immediately.
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	// Flush ensures data is written from memory buffer to file
	writer.Flush()

	// Return any error that happened during writing or flushing
	return writer.Error()
}
