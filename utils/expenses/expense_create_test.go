package expenseutils

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"

	"github.com/beego/beego/v2/server/web"
)

func setupExpenseUtilsTest(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "expenses.csv")
	if err := web.AppConfig.Set("expenses_csv_path", path); err != nil {
		t.Fatalf("set expenses csv path: %v", err)
	}
	return path
}

func TestEnsureExpensesCSV(t *testing.T) {
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
			path := setupExpenseUtilsTest(t)
			if tt.precreate {
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					t.Fatalf("mkdir: %v", err)
				}
				if err := os.WriteFile(path, []byte("custom\n"), 0644); err != nil {
					t.Fatalf("precreate file: %v", err)
				}
			}

			if err := EnsureExpensesCSV(); err != nil {
				t.Fatalf("EnsureExpensesCSV: %v", err)
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

func TestGetExpensesCSVPath(t *testing.T) {
	path := setupExpenseUtilsTest(t)
	if got := GetExpensesCSVPath(); got != path {
		t.Fatalf("path = %q, want %q", got, path)
	}
}

func TestEnsureExpensesCSVError(t *testing.T) {
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
			if err := web.AppConfig.Set("expenses_csv_path", filepath.Join(parent, "expenses.csv")); err != nil {
				t.Fatalf("set expenses csv path: %v", err)
			}
			if err := EnsureExpensesCSV(); err == nil {
				t.Fatalf("EnsureExpensesCSV error = nil, want error")
			}
		})
	}
}
