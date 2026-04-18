//store all the separate helper functions here

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
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

func (q *Queries) UpdateBook(ctx context.Context, id uuid.UUID, arg Parameters) (Book, error) {

	row := q.db.QueryRowContext(ctx, updateBook, arg.Title, ToNullString(arg.Isbn), id)

	var i Book
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Isbn,
	)
	return i, err
}
