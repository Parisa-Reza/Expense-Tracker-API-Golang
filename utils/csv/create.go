package csvutils

import (
	"encoding/csv"
	"os"
	"path/filepath"

	beego "github.com/beego/beego/v2/server/web"
)

// File path constant
const UsersCSVPath = "data/users.csv"

// CSV header row
var usersCSVHeader = []string{"id", "name", "email", "password", "created_at"}

// GetUsersCSVPath returns the configured users CSV path.
func GetUsersCSVPath() string {
	return beego.AppConfig.DefaultString("users_csv_path", UsersCSVPath)
}

// EnsureUsersCSV creates the users CSV file with a header when it does not exist.
func EnsureUsersCSV() error {
	usersCSVPath := GetUsersCSVPath()

	// Create folder if missing

	// here 0755 means: read/write permissions for owner, and read permissions for group and others
	if err := os.MkdirAll(filepath.Dir(usersCSVPath), 0755); err != nil {
		return err
	}

	// Check if file already exists
	if _, err := os.Stat(usersCSVPath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	// Create the file
	file, err := os.Create(usersCSVPath)
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
