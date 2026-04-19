package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/knibirdgautam/library/internal/database"
)

// MockDB is a manual mock of the DBQueries interface.
type MockDB struct {
	GetBookFunc func(ctx context.Context, id uuid.UUID) (database.GetBookRow, error)
	// Add other functions as needed for tests
}

func (m *MockDB) CreateAuthor(ctx context.Context, arg database.CreateAuthorParams) (database.Author, error) {
	return database.Author{}, nil
}

func (m *MockDB) CreateBook(ctx context.Context, arg database.CreateBookParams) (database.Book, error) {
	return database.Book{}, nil
}

func (m *MockDB) DeleteBook(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *MockDB) GetAuthor(ctx context.Context, name string) (database.Author, error) {
	return database.Author{}, nil
}

func (m *MockDB) LinkBookAuthor(ctx context.Context, arg database.LinkBookAuthorParams) (database.BookAuthor, error) {
	return database.BookAuthor{}, nil
}

func (m *MockDB) UnlinkBook(ctx context.Context, bookID uuid.UUID) error {
	return nil
}

func (m *MockDB) GetAuthorBooks(ctx context.Context, id uuid.UUID) ([]database.GetAuthorBooksRow, error) {
	return nil, nil
}

func (m *MockDB) GetBook(ctx context.Context, id uuid.UUID) (database.GetBookRow, error) {
	if m.GetBookFunc != nil {
		return m.GetBookFunc(ctx, id)
	}
	return database.GetBookRow{}, nil
}

func (m *MockDB) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	return database.User{}, nil
}

func (m *MockDB) UpdateBook(ctx context.Context, id uuid.UUID, arg database.Parameters) (database.Book, error) {
	return database.Book{}, nil
}
