package userutils

import (
	"testing"

	"github.com/beego/beego/v2/server/web"
)

func TestAppendUserCSV(t *testing.T) {
	tests := []struct {
		name    string
		record  []string
		wantLen int
	}{
		{
			name:    "append valid user",
			record:  []string{"1", "Parisa", "parisa@example.com", "secret1", "2025-06-10T00:00:00Z"},
			wantLen: 1,
		},
		{
			name:    "append another user",
			record:  []string{"2", "Reza", "reza@example.com", "secret1", "2025-06-10T00:00:00Z"},
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupUserUtilsTest(t)
			if err := AppendUserCSV(tt.record); err != nil {
				t.Fatalf("AppendUserCSV: %v", err)
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

func TestAppendUserCSVError(t *testing.T) {
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
			record := []string{"1", "Parisa", "parisa@example.com", "secret1", "2025-06-10T00:00:00Z"}
			if err := AppendUserCSV(record); err == nil {
				t.Fatalf("AppendUserCSV error = nil, want error")
			}
		})
	}
}
