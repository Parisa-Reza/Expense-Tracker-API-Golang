package csvutils

import (
	"encoding/csv"
	"os"
	"path/filepath"
)

// File path constant
const UsersCSVPath = "data/users.csv"

// CSV header row
var usersCSVHeader = []string{"id", "name", "email", "password", "created_at"}

func EnsureUsersCSV() error {

	// Create folder if missing

	// here 0755 means: read/write permissions for owner, and read permissions for group and others
	if err := os.MkdirAll(filepath.Dir(UsersCSVPath), 0755); err != nil {
		return err
	}

	// Check if file already exists
	if _, err := os.Stat(UsersCSVPath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}


	// Create the file
	file, err := os.Create(UsersCSVPath)
	if err != nil {
		return err
	}

	// Close the file when function ends
	defer file.Close()


	// : Write header row
	writer := csv.NewWriter(file)
	if err := writer.Write(usersCSVHeader); err != nil {
		return err
	}

	// Flush() forces buffered (temporary) CSV data to be permanently written into the file.
	writer.Flush()


	// If anything went wrong while writing, return the error. Otherwise returns nil
	return writer.Error()
}
