package models

import (
	"errors"
	"testing"
)

func TestUserStorageOperations(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "create and get by email",
			run: func(t *testing.T) {
				user := User{Name: "Parisa", Email: "parisa@example.com", Password: "secret1"}
				if err := CreateUser(&user); err != nil {
					t.Fatalf("create user: %v", err)
				}
				got, err := GetUserByEmail("PARISA@example.com")
				if err != nil {
					t.Fatalf("get by email: %v", err)
				}
				if got.ID != user.ID {
					t.Fatalf("id = %d, want %d", got.ID, user.ID)
				}
			},
		},
		{
			name: "get by id",
			run: func(t *testing.T) {
				user := User{Name: "Reza", Email: "reza@example.com", Password: "secret1"}
				if err := CreateUser(&user); err != nil {
					t.Fatalf("create user: %v", err)
				}
				got, err := GetUserByID(user.ID)
				if err != nil {
					t.Fatalf("get by id: %v", err)
				}
				if got.Email != user.Email {
					t.Fatalf("email = %q, want %q", got.Email, user.Email)
				}
			},
		},
		{
			name: "valid user id",
			run: func(t *testing.T) {
				user := User{Name: "Valid", Email: "valid@example.com", Password: "secret1"}
				if err := CreateUser(&user); err != nil {
					t.Fatalf("create user: %v", err)
				}
				if !IsValidUserID(user.ID) {
					t.Fatalf("IsValidUserID(%d) = false, want true", user.ID)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupModelStorage(t)
			tt.run(t)
		})
	}
}

func TestUserFailures(t *testing.T) {
	tests := []struct {
		name    string
		run     func() error
		wantErr error
	}{
		{name: "missing email", run: func() error { _, err := GetUserByEmail("missing@example.com"); return err }, wantErr: ErrUserNotFound},
		{name: "missing id", run: func() error { _, err := GetUserByID(42); return err }, wantErr: ErrUserNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupModelStorage(t)
			if err := tt.run(); !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetNextID(t *testing.T) {
	setupModelStorage(t)
	if got := GetNextID(); got != 1 {
		t.Fatalf("empty next id = %d, want 1", got)
	}
	if err := CreateUser(&User{Name: "Parisa", Email: "parisa@example.com", Password: "secret1"}); err != nil {
		t.Fatalf("create user: %v", err)
	}
	if got := GetNextID(); got != 2 {
		t.Fatalf("next id = %d, want 2", got)
	}
}
