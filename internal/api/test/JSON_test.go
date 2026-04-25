package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondWithJSON(t *testing.T) {
	w := httptest.NewRecorder()
	payload := map[string]string{"message": "hello"}
	statusCode := http.StatusOK

	RespondWithJSON(w, statusCode, payload)

	if w.Code != statusCode {
		t.Errorf("Expected status code %d, got %d", statusCode, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "hello" {
		t.Errorf("Expected message 'hello', got '%s'", response["message"])
	}
}

func TestRespondWithError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		msg        string
		expected   string
	}{
		{
			name:       "Client Error",
			statusCode: http.StatusBadRequest,
			msg:        "Bad Request",
			expected:   "Bad Request",
		},
		{
			name:       "Server Error Masking",
			statusCode: http.StatusInternalServerError,
			msg:        "Database Exploded",
			expected:   "Internal Server Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			RespondWithError(w, tt.statusCode, tt.msg)

			if w.Code != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, w.Code)
			}

			var response map[string]string
			json.Unmarshal(w.Body.Bytes(), &response)

			if response["Error"] != tt.expected {
				t.Errorf("Expected error message '%s', got '%s'", tt.expected, response["Error"])
			}
		})
	}
}
