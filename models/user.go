package models

import (
	"errors"
	"strconv"
	"strings"
	"time"

	csvutils "expense-tracker-api/utils/csv"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	ErrUserNotFound = errors.New("user not found")
)

func GetAllUsers() ([]User, error) {
	records, err := csvutils.ReadUsersCSV()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(records))
	for _, record := range records {
		if len(record) < 5 {
			continue
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			continue
		}

		createdAt, err := time.Parse(time.RFC3339, record[4])
		if err != nil {
			continue
		}

		users = append(users, User{
			ID:        id,
			Name:      record[1],
			Email:     record[2],
			Password:  record[3],
			CreatedAt: createdAt,
		})
	}

	return users, nil
}

func GetUserByEmail(email string) (*User, error) {
	users, err := GetAllUsers()
	if err != nil {
		return nil, err
	}

	for i := range users {
		if strings.EqualFold(users[i].Email, strings.TrimSpace(email)) {
			return &users[i], nil
		}
	}

	return nil, ErrUserNotFound
}

func GetUserByID(id int) (*User, error) {
	users, err := GetAllUsers()
	if err != nil {
		return nil, err
	}

	for i := range users {
		if users[i].ID == id {
			return &users[i], nil
		}
	}

	return nil, ErrUserNotFound
}

func CreateUser(user *User) error {
	user.ID = GetNextID()
	user.CreatedAt = time.Now().UTC()

	return csvutils.AppendUserCSV([]string{
		strconv.Itoa(user.ID),
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt.Format(time.RFC3339),
	})
}

func GetNextID() int {
	users, err := GetAllUsers()
	if err != nil || len(users) == 0 {
		return 1
	}

	maxID := 0
	// here we loop through all users to find the maximum ID and return maxID + 1 as the next ID
	for _, user := range users {
		if user.ID > maxID {
			maxID = user.ID
		}
	}

	return maxID + 1
}

func IsValidUserID(id int) bool {
	_, err := GetUserByID(id)
	return err == nil
}
