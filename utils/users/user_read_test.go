package userutils

import (
	"testing"

	"github.com/beego/beego/v2/server/web"
)

func TestReadUsersCSV(t *testing.T) {
	tests := []struct {
		name    string
		records [][]string
		wantLen int
	}{
		{name: "empty csv", wantLen: 0},
		{
			name: "with records",
			records: [][]string{
				{"1", "Parisa", "parisa@example.com", "secret1", "2025-06-10T00:00:00Z"},
				{"2", "Reza", "reza@example.com", "secret1", "2025-06-10T00:00:00Z"},
			},
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupUserUtilsTest(t)
			for _, record := range tt.records {
				if err := AppendUserCSV(record); err != nil {
					t.Fatalf("AppendUserCSV: %v", err)
				}
			}
			got, err := ReadUsersCSV()
			if err != nil {
				t.Fatalf("ReadUsersCSV: %v", err)
			}
			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestReadUsersCSVError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "path is directory"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := web.AppConfig.Set("users_csv_path", t.TempDir()); err != nil {
				t.Fatalf("set users csv path: %v", err)
			}
			if _, err := ReadUsersCSV(); err == nil {
				t.Fatalf("ReadUsersCSV error = nil, want error")
			}
		})
	}
}
