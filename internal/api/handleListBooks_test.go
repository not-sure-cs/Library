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

func TestHandleListOfAuthorBooks(t *testing.T) {
	authorID := uuid.New()

	tests := []struct {
		name           string
		id             string
		mockList       func(ctx context.Context, id uuid.UUID) ([]database.GetAuthorBooksRow, error)
		expectedStatus int
	}{
		{
			name: "Success",
			id:   authorID.String(),
			mockList: func(ctx context.Context, id uuid.UUID) ([]database.GetAuthorBooksRow, error) {
				return []database.GetAuthorBooksRow{
					{Name: "Book 1", AuthorID: authorID},
					{Name: "Book 2", AuthorID: authorID},
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error - Invalid UUID",
			id:             "bad-uuid",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Error - DB Failure",
			id:             authorID.String(),
			mockList:       func(ctx context.Context, id uuid.UUID) ([]database.GetAuthorBooksRow, error) { return nil, errors.New("db error") },
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				GetAuthorBooksFunc: tt.mockList,
			}

			mux := http.NewServeMux()
			mux.HandleFunc("GET /author/{id}/books", HandleListOfAuthorBooks(mockDB))

			req := httptest.NewRequest("GET", "/author/"+tt.id+"/books", nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var resp []database.GetAuthorBooksRow
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if len(resp) != 2 {
					t.Errorf("Expected 2 books, got %d", len(resp))
				}
			}
		})
	}
}
