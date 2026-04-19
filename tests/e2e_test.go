package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/knibirdgautam/library/internal/api"
	"github.com/knibirdgautam/library/internal/database"
)

// MockDB for E2E demonstration
type MockDB struct {
	GetBookFunc func(ctx context.Context, id uuid.UUID) (database.GetBookRow, error)
}

func (m *MockDB) CreateAuthor(ctx context.Context, arg database.CreateAuthorParams) (database.Author, error) { return database.Author{}, nil }
func (m *MockDB) CreateBook(ctx context.Context, arg database.CreateBookParams) (database.Book, error) { return database.Book{}, nil }
func (m *MockDB) DeleteBook(ctx context.Context, id uuid.UUID) error { return nil }
func (m *MockDB) GetAuthor(ctx context.Context, name string) (database.Author, error) { return database.Author{}, nil }
func (m *MockDB) LinkBookAuthor(ctx context.Context, arg database.LinkBookAuthorParams) (database.BookAuthor, error) { return database.BookAuthor{}, nil }
func (m *MockDB) UnlinkBook(ctx context.Context, bookID uuid.UUID) error { return nil }
func (m *MockDB) GetAuthorBooks(ctx context.Context, id uuid.UUID) ([]database.GetAuthorBooksRow, error) { return nil, nil }
func (m *MockDB) GetBook(ctx context.Context, id uuid.UUID) (database.GetBookRow, error) {
	if m.GetBookFunc != nil {
		return m.GetBookFunc(ctx, id)
	}
	return database.GetBookRow{}, nil
}
func (m *MockDB) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) { return database.User{}, nil }
func (m *MockDB) UpdateBook(ctx context.Context, id uuid.UUID, arg database.Parameters) (database.Book, error) { return database.Book{}, nil }

func TestServerE2E(t *testing.T) {
	// 1. Setup
	mockDB := &MockDB{
		GetBookFunc: func(ctx context.Context, id uuid.UUID) (database.GetBookRow, error) {
			return database.GetBookRow{
				BookID: id,
				Name:   "E2E Test Book",
			}, nil
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /book/{id}", api.HandleGetBooks(mockDB))
	
	// Apply middleware just like in main.go
	handler := api.JSONMiddleware(mux)
	
	server := httptest.NewServer(handler)
	defer server.Close()

	// 2. Execution
	bookID := uuid.New()
	resp, err := http.Get(server.URL + "/book/" + bookID.String())
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// 3. Assertion
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected JSON content type, got %v", contentType)
	}

	var book database.GetBookRow
	if err := json.NewDecoder(resp.Body).Decode(&book); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if book.Name != "E2E Test Book" {
		t.Errorf("Expected 'E2E Test Book', got '%s'", book.Name)
	}
}

func TestHealthCheckE2E(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /status", api.HandleStatus(time.Now())) // Note: time dependency usually injected
	
	server := httptest.NewServer(api.JSONMiddleware(mux))
	defer server.Close()

	resp, err := http.Get(server.URL + "/status")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}
