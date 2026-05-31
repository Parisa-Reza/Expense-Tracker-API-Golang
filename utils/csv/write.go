package csvutils

import (
	"encoding/csv"
	"os"
)

// AppendUserCSV adds a new user record (row) into the users CSV file
func AppendUserCSV(record []string) error {

	// Make sure the CSV file exists (create it if not)
	if err := EnsureUsersCSV(); err != nil {
		return err
	}

	// Open the file in append mode so we can add new data at the end . 0644 means: read/write permissions for owner, and read permissions for group and others
	file, err := os.OpenFile(GetUsersCSVPath(), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Close file when function finishes
	defer file.Close()

	// Create a CSV writer to write data into the file
	writer := csv.NewWriter(file)

	// Write one row (user record) into CSV

	if err := writer.Write(record); err != nil {
		return err
	}

	// Flush ensures data is written from memory buffer to file
	writer.Flush()

	// Return any error that happened during writing or flushing
	return writer.Error()
}
