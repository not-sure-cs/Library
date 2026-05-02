package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

// DBQueries is an interface that wraps all the database methods.
// This allows us to mock the database during testing.
type DBQueries interface {
	CreateAuthor(ctx context.Context, arg CreateAuthorParams) (Author, error)
	CreateBook(ctx context.Context, arg CreateBookParams) (Book, error)
	DeleteBook(ctx context.Context, id uuid.UUID) error
	GetAuthor(ctx context.Context, name string) (Author, error)
	LinkBookAuthor(ctx context.Context, arg LinkBookAuthorParams) (BookAuthor, error)
	UnlinkBook(ctx context.Context, bookID uuid.UUID) error
	GetAuthorBooks(ctx context.Context, id uuid.UUID) ([]GetAuthorBooksRow, error)
	GetBook(ctx context.Context, id uuid.UUID) (GetBookRow, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	UpdateBook(ctx context.Context, id uuid.UUID, arg Parameters) (Book, error)
	CountBook(ctx context.Context) (int64, error)
	LinkHash(ctx context.Context, arg LinkHashParams) error
	GetUser(ctx context.Context, email sql.NullString) (User, error)
	GetPassHash(ctx context.Context, email sql.NullString) (string, error)
}
