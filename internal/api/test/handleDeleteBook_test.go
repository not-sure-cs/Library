package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestHandleDeleteBook(t *testing.T) {
	bookID := uuid.New()

	tests := []struct {
		name           string
		id             string
		mockUnlink     func(ctx context.Context, id uuid.UUID) error
		mockDelete     func(ctx context.Context, id uuid.UUID) error
		expectedStatus int
	}{
		{
			name:           "Success",
			id:             bookID.String(),
			mockUnlink:     func(ctx context.Context, id uuid.UUID) error { return nil },
			mockDelete:     func(ctx context.Context, id uuid.UUID) error { return nil },
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error - Invalid UUID",
			id:             "bad-uuid",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Error - Unlink Failure",
			id:             bookID.String(),
			mockUnlink:     func(ctx context.Context, id uuid.UUID) error { return errors.New("unlink failed") },
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Error - Delete Failure",
			id:   bookID.String(),
			mockUnlink: func(ctx context.Context, id uuid.UUID) error { return nil },
			mockDelete: func(ctx context.Context, id uuid.UUID) error { return errors.New("delete failed") },
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				UnlinkBookFunc: tt.mockUnlink,
				DeleteBookFunc: tt.mockDelete,
			}

			mux := http.NewServeMux()
			mux.HandleFunc("DELETE /book/{id}", HandleDeleteBook(mockDB))

			req := httptest.NewRequest("DELETE", "/book/"+tt.id, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
