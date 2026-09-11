//store all the separate helper functions here

package database

import (
	"context"
	"database/sql"
	"log"
	"mime/multipart"
	"path/filepath"
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

func ToNullInt16(s int16) sql.NullInt16 {
	var r sql.NullInt16

	r.Int16 = s
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
        COALESCE(a.name, '') as author_name
    FROM updated_book ub
    LEFT JOIN book_authors ba ON ub.id = ba.book_id
    LEFT JOIN authors a ON ba.author_id = a.id
    LIMIT 1;
`

func (q *Queries) UpdateBook(ctx context.Context, id uuid.UUID, arg Parameters) (UserBook, error) {
	if arg.Author != "" {
		author, err := q.GetAuthor(ctx, arg.Author)
		if err != nil {
			author, err = q.CreateAuthor(ctx, CreateAuthorParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Name:      arg.Author,
			})
			if err != nil {
				return UserBook{}, err
			}
		}
		_ = q.UnlinkBook(ctx, id)
		_, err = q.LinkBookAuthor(ctx, LinkBookAuthorParams{
			BookID:   id,
			AuthorID: author.ID,
		})
		if err != nil {
			return UserBook{}, err
		}
	}

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
	return "Type-Book-" + uuid.New().String()
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
	ext := filepath.Ext(handler.Filename)
	name := GenerateFileName(total, time.Now()) + ext
	key := path + name

	contentType := handler.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	err := store.UploadFile(r, secret.Bucket, key, contentType, file)
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

func PasswordHash(pass []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Bcrypt password generation error: %v", err)
		return "", err
	}
	return string(hash), nil
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
