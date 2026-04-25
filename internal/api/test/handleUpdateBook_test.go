package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/knibirdgautam/library/internal/database"
)

func TestHandleUpdateBooks(t *testing.T) {
	bookID := uuid.New()

	tests := []struct {
		name           string
		id             string
		body           interface{}
		mockUpdate     func(ctx context.Context, id uuid.UUID, arg database.Parameters) (database.Book, error)
		expectedStatus int
	}{
		{
			name: "Success",
			id:   bookID.String(),
			body: database.Parameters{
				Title: "Updated Title",
				Isbn:  "654321",
			},
			mockUpdate: func(ctx context.Context, id uuid.UUID, arg database.Parameters) (database.Book, error) {
				return database.Book{ID: id, Name: arg.Title}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error - Invalid UUID",
			id:             "bad-uuid",
			body:           database.Parameters{Title: "Title"},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Error - Invalid Body",
			id:             bookID.String(),
			body:           "not-a-struct",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Error - DB Failure",
			id:   bookID.String(),
			body: database.Parameters{Title: "Title"},
			mockUpdate: func(ctx context.Context, id uuid.UUID, arg database.Parameters) (database.Book, error) {
				return database.Book{}, errors.New("db error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				UpdateBookFunc: tt.mockUpdate,
			}

			mux := http.NewServeMux()
			mux.HandleFunc("PUT /book/{id}", HandleUpdateBooks(mockDB))

			var bodyReader *bytes.Reader
			if s, ok := tt.body.(string); ok {
				bodyReader = bytes.NewReader([]byte(s))
			} else {
				jsonBytes, _ := json.Marshal(tt.body)
				bodyReader = bytes.NewReader(jsonBytes)
			}

			req := httptest.NewRequest("PUT", "/book/"+tt.id, bodyReader)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
