package expenseutils

import (
	"testing"

	"github.com/beego/beego/v2/server/web"
)

func TestReadExpensesCSV(t *testing.T) {
	tests := []struct {
		name    string
		records [][]string
		wantLen int
	}{
		{name: "empty csv", wantLen: 0},
		{
			name: "with records",
			records: [][]string{
				{"1", "1", "Lunch", "350.50", "Food", "", "2025-06-10", "2025-06-10T00:00:00Z"},
				{"2", "1", "Bus", "50.00", "Transport", "", "2025-06-11", "2025-06-10T00:00:00Z"},
			},
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupExpenseUtilsTest(t)
			if len(tt.records) > 0 {
				if err := WriteExpensesCSV(tt.records); err != nil {
					t.Fatalf("WriteExpensesCSV: %v", err)
				}
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

func TestReadExpensesCSVError(t *testing.T) {
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
			if _, err := ReadExpensesCSV(); err == nil {
				t.Fatalf("ReadExpensesCSV error = nil, want error")
			}
		})
	}
}
