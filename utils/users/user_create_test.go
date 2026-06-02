package userutils

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"

	"github.com/beego/beego/v2/server/web"
)

func setupUserUtilsTest(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "users.csv")
	if err := web.AppConfig.Set("users_csv_path", path); err != nil {
		t.Fatalf("set users csv path: %v", err)
	}
	return path
}

func TestEnsureUsersCSV(t *testing.T) {
	tests := []struct {
		name       string
		precreate  bool
		wantHeader bool
	}{
		{name: "creates missing csv", wantHeader: true},
		{name: "keeps existing csv", precreate: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := setupUserUtilsTest(t)
			if tt.precreate {
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					t.Fatalf("mkdir: %v", err)
				}
				if err := os.WriteFile(path, []byte("custom\n"), 0644); err != nil {
					t.Fatalf("precreate file: %v", err)
				}
			}

			if err := EnsureUsersCSV(); err != nil {
				t.Fatalf("EnsureUsersCSV: %v", err)
			}
			if _, err := os.Stat(path); err != nil {
				t.Fatalf("stat csv: %v", err)
			}

			if tt.wantHeader {
				file, err := os.Open(path)
				if err != nil {
					t.Fatalf("open csv: %v", err)
				}
				defer file.Close()
				records, err := csv.NewReader(file).ReadAll()
				if err != nil {
					t.Fatalf("read csv: %v", err)
				}
				if len(records) != 1 || records[0][0] != "id" {
					t.Fatalf("header = %#v", records)
				}
			}
		})
	}
}

func TestGetUsersCSVPath(t *testing.T) {
	path := setupUserUtilsTest(t)
	if got := GetUsersCSVPath(); got != path {
		t.Fatalf("path = %q, want %q", got, path)
	}
}

func TestEnsureUsersCSVError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "parent path is file"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parent := filepath.Join(t.TempDir(), "parent-file")
			if err := os.WriteFile(parent, []byte("not a dir"), 0644); err != nil {
				t.Fatalf("write parent file: %v", err)
			}
			if err := web.AppConfig.Set("users_csv_path", filepath.Join(parent, "users.csv")); err != nil {
				t.Fatalf("set users csv path: %v", err)
			}
			if err := EnsureUsersCSV(); err == nil {
				t.Fatalf("EnsureUsersCSV error = nil, want error")
			}
		})
	}
}
