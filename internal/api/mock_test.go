package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/knibirdgautam/library/internal/database"
)

// MockDB is a manual mock of the DBQueries interface.
type MockDB struct {
	CreateAuthorFunc   func(ctx context.Context, arg database.CreateAuthorParams) (database.Author, error)
	CreateBookFunc     func(ctx context.Context, arg database.CreateBookParams) (database.Book, error)
	DeleteBookFunc     func(ctx context.Context, id uuid.UUID) error
	GetAuthorFunc       func(ctx context.Context, name string) (database.Author, error)
	LinkBookAuthorFunc func(ctx context.Context, arg database.LinkBookAuthorParams) (database.BookAuthor, error)
	UnlinkBookFunc     func(ctx context.Context, bookID uuid.UUID) error
	GetAuthorBooksFunc func(ctx context.Context, id uuid.UUID) ([]database.GetAuthorBooksRow, error)
	GetBookFunc        func(ctx context.Context, id uuid.UUID) (database.GetBookRow, error)
	CreateUserFunc     func(ctx context.Context, arg database.CreateUserParams) (database.User, error)
	UpdateBookFunc     func(ctx context.Context, id uuid.UUID, arg database.Parameters) (database.Book, error)
}

func (m *MockDB) CreateAuthor(ctx context.Context, arg database.CreateAuthorParams) (database.Author, error) {
	if m.CreateAuthorFunc != nil {
		return m.CreateAuthorFunc(ctx, arg)
	}
	return database.Author{}, nil
}

func (m *MockDB) CreateBook(ctx context.Context, arg database.CreateBookParams) (database.Book, error) {
	if m.CreateBookFunc != nil {
		return m.CreateBookFunc(ctx, arg)
	}
	return database.Book{}, nil
}

func (m *MockDB) DeleteBook(ctx context.Context, id uuid.UUID) error {
	if m.DeleteBookFunc != nil {
		return m.DeleteBookFunc(ctx, id)
	}
	return nil
}

func (m *MockDB) GetAuthor(ctx context.Context, name string) (database.Author, error) {
	if m.GetAuthorFunc != nil {
		return m.GetAuthorFunc(ctx, name)
	}
	return database.Author{}, nil
}

func (m *MockDB) LinkBookAuthor(ctx context.Context, arg database.LinkBookAuthorParams) (database.BookAuthor, error) {
	if m.LinkBookAuthorFunc != nil {
		return m.LinkBookAuthorFunc(ctx, arg)
	}
	return database.BookAuthor{}, nil
}

func (m *MockDB) UnlinkBook(ctx context.Context, bookID uuid.UUID) error {
	if m.UnlinkBookFunc != nil {
		return m.UnlinkBookFunc(ctx, bookID)
	}
	return nil
}

func (m *MockDB) GetAuthorBooks(ctx context.Context, id uuid.UUID) ([]database.GetAuthorBooksRow, error) {
	if m.GetAuthorBooksFunc != nil {
		return m.GetAuthorBooksFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockDB) GetBook(ctx context.Context, id uuid.UUID) (database.GetBookRow, error) {
	if m.GetBookFunc != nil {
		return m.GetBookFunc(ctx, id)
	}
	return database.GetBookRow{}, nil
}

func (m *MockDB) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, arg)
	}
	return database.User{}, nil
}

func (m *MockDB) UpdateBook(ctx context.Context, id uuid.UUID, arg database.Parameters) (database.Book, error) {
	if m.UpdateBookFunc != nil {
		return m.UpdateBookFunc(ctx, id, arg)
	}
	return database.Book{}, nil
}
