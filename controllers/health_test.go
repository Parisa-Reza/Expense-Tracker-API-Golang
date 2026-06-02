package controllers

import (
	"net/http"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	tests := []struct {
		name       string
		wantStatus int
		wantMsg    string
	}{
		{name: "server running", wantStatus: http.StatusOK, wantMsg: "Server is running"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupControllerTest(t)
			response := performRequest(http.MethodGet, "/api/v1/health", "", "")
			got := decodeAPIResponse(t, response)

			if response.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, tt.wantStatus)
			}
			if got.Message != tt.wantMsg {
				t.Fatalf("message = %q, want %q", got.Message, tt.wantMsg)
			}
		})
	}
}
