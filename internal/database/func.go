//store all the separate helper functions here

package database

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

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

const path = "Library/Public/Books/"

func GenerateFileName(id int64, now time.Time) string {

	t := now.Unix()

	fileName := "Type-Book" + strconv.Itoa(int(id)) + strconv.Itoa(int(t))
	return fileName

}

const countBook = `
	SELECT COUNT(*) FROM books ;
`

func (q *Queries) CountBook(ctx context.Context) (int64, error) {

	var store int64
	row := q.db.QueryRowContext(ctx, countBook)

	err := row.Scan(&store)
	return store, err

}

func SaveFile(total int64, file multipart.File, handler *multipart.FileHeader) (string, error) {

	defer file.Close()

	tempFolderPath := fmt.Sprintf("./%s", path)
	tempFileName := fmt.Sprintf("upload-%s-*.%s", GenerateFileName(total, time.Now()), filepath.Ext(handler.Filename))

	tempFile, err := os.CreateTemp(tempFolderPath, tempFileName)
	if err != nil {
		return "", err
	}

	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {

		return "", err
	}

	tempFile.Write(fileBytes)
	return tempFile.Name(), nil
}
