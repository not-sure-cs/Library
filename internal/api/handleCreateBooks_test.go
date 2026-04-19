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

func TestHandleCreateBooks(t *testing.T) {
	authorID := uuid.New()
	bookID := uuid.New()

	tests := []struct {
		name           string
		body           map[string]string
		mockGetAuthor  func(ctx context.Context, name string) (database.Author, error)
		mockCreateAuth func(ctx context.Context, arg database.CreateAuthorParams) (database.Author, error)
		mockCreateBook func(ctx context.Context, arg database.CreateBookParams) (database.Book, error)
		mockLink       func(ctx context.Context, arg database.LinkBookAuthorParams) (database.BookAuthor, error)
		expectedStatus int
	}{
		{
			name: "Success - Author Exists",
			body: map[string]string{
				"title":  "Test Book",
				"isbn":   "123456",
				"author": "Test Author",
			},
			mockGetAuthor: func(ctx context.Context, name string) (database.Author, error) {
				return database.Author{ID: authorID, Name: name}, nil
			},
			mockCreateBook: func(ctx context.Context, arg database.CreateBookParams) (database.Book, error) {
				return database.Book{ID: bookID, Name: arg.Name}, nil
			},
			mockLink: func(ctx context.Context, arg database.LinkBookAuthorParams) (database.BookAuthor, error) {
				return database.BookAuthor{BookID: bookID, AuthorID: authorID}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Success - Author Created",
			body: map[string]string{
				"title":  "Test Book",
				"isbn":   "123456",
				"author": "New Author",
			},
			mockGetAuthor: func(ctx context.Context, name string) (database.Author, error) {
				return database.Author{}, errors.New("not found")
			},
			mockCreateAuth: func(ctx context.Context, arg database.CreateAuthorParams) (database.Author, error) {
				return database.Author{ID: authorID, Name: arg.Name}, nil
			},
			mockCreateBook: func(ctx context.Context, arg database.CreateBookParams) (database.Book, error) {
				return database.Book{ID: bookID, Name: arg.Name}, nil
			},
			mockLink: func(ctx context.Context, arg database.LinkBookAuthorParams) (database.BookAuthor, error) {
				return database.BookAuthor{BookID: bookID, AuthorID: authorID}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error - Invalid Body",
			body:           nil, // Will result in empty body or decoder error
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Error - DB Failure",
			body: map[string]string{
				"title":  "Test Book",
				"isbn":   "123456",
				"author": "Test Author",
			},
			mockGetAuthor: func(ctx context.Context, name string) (database.Author, error) {
				return database.Author{ID: authorID, Name: name}, nil
			},
			mockCreateBook: func(ctx context.Context, arg database.CreateBookParams) (database.Book, error) {
				return database.Book{}, errors.New("db error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				GetAuthorFunc:      tt.mockGetAuthor,
				CreateAuthorFunc:   tt.mockCreateAuth,
				CreateBookFunc:     tt.mockCreateBook,
				LinkBookAuthorFunc: tt.mockLink,
			}

			handler := HandleCreateBooks(mockDB)

			var bodyReader *bytes.Reader
			if tt.body != nil {
				jsonBytes, _ := json.Marshal(tt.body)
				bodyReader = bytes.NewReader(jsonBytes)
			} else {
				bodyReader = bytes.NewReader([]byte("{invalid json}"))
			}

			req := httptest.NewRequest("POST", "/book", bodyReader)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var resp map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				book := resp["book"].(map[string]interface{})
				if book["Name"] != tt.body["title"] {
					t.Errorf("Expected book title %s, got %s", tt.body["title"], book["Name"])
				}
			}
		})
	}
}
