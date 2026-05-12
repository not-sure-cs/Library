//store all the separate helper functions here

package database

import (
	"context"
	"database/sql"
	"log"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/knibirdgautam/library/internal/storage"
	"golang.org/x/crypto/bcrypt"
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
    WITH updated_book AS (
        UPDATE books
        SET name = $1, isbn = $2, updated_at = NOW()
        WHERE id = $3
        RETURNING id, name, isbn
    )
    SELECT 
        ub.name, 
        ub.isbn, 
        a.name as author_name
    FROM updated_book ub
    JOIN book_authors ba ON ub.id = ba.book_id
    JOIN authors a ON ba.author_id = a.id;
`

func (q *Queries) UpdateBook(ctx context.Context, id uuid.UUID, arg Parameters) (UserBook, error) {

	row := q.db.QueryRowContext(ctx, updateBook, arg.Title, ToNullString(arg.Isbn), id)

	var ub UserBook
	err := row.Scan(
		&ub.BookName,
		&ub.ISBN,
		&ub.AuthorName,
	)

	return ub, err

}

const path = "Assets/Books/"

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

func SaveFile(total int64, r context.Context, secret storage.Secret, store storage.R2Store, file multipart.File, handler *multipart.FileHeader) (string, error) {

	defer file.Close()

	ext := filepath.Ext(handler.Filename)
	name := GenerateFileName(total, time.Now()) + ext
	key := path + name

	err := store.UploadFile(r, secret.Bucket, key, handler.Header.Get("contentType"), file)
	if err != nil {
		return "", err
	}

	return key, nil

}

func UnsaveFile(r context.Context, secret storage.Secret, store storage.R2Store, filename string) error {
	err := store.DeleteFile(r, secret.Bucket, filename)
	if err != nil {
		return err
	}

	return nil

}

func PasswordHash(pass []byte) string {

	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func PasswordVerify(hash string, pass []byte) bool {
	byteHash := []byte(hash)
	err := bcrypt.CompareHashAndPassword(byteHash, pass)
	if err != nil {
		log.Println(err)
		return false
	}
	return true

}
