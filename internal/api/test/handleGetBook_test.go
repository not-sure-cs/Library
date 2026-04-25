package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/knibirdgautam/library/internal/database"
)

func TestHandleGetBooks(t *testing.T) {
	bookID := uuid.New()
	
	tests := []struct {
		name           string
		id             string
		mockReturn     database.GetBookRow
		mockErr        error
		expectedStatus int
	}{
		{
			name: "Success",
			id:   bookID.String(),
			mockReturn: database.GetBookRow{
				BookID: bookID,
				Name:   "The Go Programming Language",
			},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid ID",
			id:             "not-a-uuid",
			mockReturn:     database.GetBookRow{},
			mockErr:        nil,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Not Found",
			id:             bookID.String(),
			mockReturn:     database.GetBookRow{},
			mockErr:        errors.New("not found"),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				GetBookFunc: func(ctx context.Context, id uuid.UUID) (database.GetBookRow, error) {
					return tt.mockReturn, tt.mockErr
				},
			}

			// We use a mux to handle the path parameter {id}
			mux := http.NewServeMux()
			mux.HandleFunc("GET /book/{id}", HandleGetBooks(mockDB))

			req := httptest.NewRequest("GET", "/book/"+tt.id, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var resp database.GetBookRow
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if resp.BookID != bookID {
					t.Errorf("Expected book ID %v, got %v", bookID, resp.BookID)
				}
			}
		})
	}
}
