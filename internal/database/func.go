//store all the separate helper functions here

package database

import (
	"context"
	"database/sql"

)

func ToNullString(s string) sql.NullString {
	var r sql.NullString
	if s == "" {
		r.Valid = false
		return r
	}

	r.String = s
	r.Valid = true
	return r
}

type Parameters struct {
	Title  string `json:"title"`
	Isbn   string `json:"isbn"`
	Author string `json:"author"`
}

const updateBook = `
	UPDATE books
	SET name = $1, isbn = $2, updated_at = NOW()
	WHERE id = $3
	RETURNING id, created_at, updated_at, name, isbn
`

func (q *Queries) UpdateBook(ctx context.Context, book *GetBookRow, arg Parameters) (*GetBookRow, error) {

	row := q.db.QueryRowContext(ctx, updateBook, arg.Title, arg.Author, arg.Isbn,)

	err := row.Scan(
		&book.Name,
		&book.Name_2,
		&book.BookID,
		&book.CreatedAt,
		&book.UpdatedAt,
		&book.Name,
		&book.Isbn,
	)
	return book, err
}
