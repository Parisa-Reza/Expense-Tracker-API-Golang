// package csvutils

// import (
// 	"encoding/csv"
// 	"os"
// )

// func ReadUsersCSV() ([][]string, error) {
// 	if err := EnsureUsersCSV(); err != nil {
// 		return nil, err
// 	}

// 	file, err := os.Open(UsersCSVPath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	records, err := csv.NewReader(file).ReadAll()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(records) <= 1 {
// 		return [][]string{}, nil
// 	}

// 	return records[1:], nil
// }


package csvutils

import (
	"encoding/csv"
	"os"
)

// ReadUsersCSV reads all user records from the users CSV file. It returns a 2D slice of strings where each inner slice is one user row.
func ReadUsersCSV() ([][]string, error) {

	// Ensure the CSV file exists 
	if err := EnsureUsersCSV(); err != nil {
		return nil, err
	}

	// Open the CSV file for reading
	file, err := os.Open(UsersCSVPath)
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