package expenseutils

import (
	"path/filepath"
	"testing"

	"github.com/beego/beego/v2/server/web"
)

func TestAppendExpenseCSV(t *testing.T) {
	tests := []struct {
		name    string
		record  []string
		wantLen int
	}{
		{
			name:    "append valid expense",
			record:  []string{"1", "1", "Lunch", "350.50", "Food", "", "2025-06-10", "2025-06-10T00:00:00Z"},
			wantLen: 1,
		},
		{
			name:    "append second expense",
			record:  []string{"2", "1", "Bus", "50.00", "Transport", "", "2025-06-11", "2025-06-10T00:00:00Z"},
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupExpenseUtilsTest(t)
			if err := AppendExpenseCSV(tt.record); err != nil {
				t.Fatalf("AppendExpenseCSV: %v", err)
			}
			got, err := ReadExpensesCSV()
			if err != nil {
				t.Fatalf("ReadExpensesCSV: %v", err)
			}
			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestAppendExpenseCSVError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "path is directory"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := web.AppConfig.Set("expenses_csv_path", t.TempDir()); err != nil {
				t.Fatalf("set expenses csv path: %v", err)
			}
			record := []string{"1", "1", "Lunch", "350.50", "Food", "", "2025-06-10", "2025-06-10T00:00:00Z"}
			if err := AppendExpenseCSV(record); err == nil {
				t.Fatalf("AppendExpenseCSV error = nil, want error")
			}
			if filepath.Base(GetExpensesCSVPath()) == "" {
				t.Fatalf("expected configured path")
			}
		})
	}
}
