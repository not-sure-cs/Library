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

func ToNullInt32(s int32) sql.NullInt32 {
	var r sql.NullInt32
	if s == 0 {
		r.Valid = false
		return r
	}

	r.Int32 = s
	r.Valid = true
	return r
}

type Parameters struct {
	Title     string `json:"title"`
	Isbn      string `json:"isbn"`
	Author    string `json:"author"`
	Subject   string `json:"subject,omitempty"`
	Producer  string `json:"producer,omitempty"`
	PageCount int32  `json:"page_count,omitempty"`
}

const updateBook = `
    WITH updated_book AS (
        UPDATE books
        SET name = COALESCE(NULLIF($1, ''), name),
            isbn = CASE WHEN $2::text = '' THEN isbn ELSE $2 END,
            subject = COALESCE(NULLIF($3, ''), subject),
            producer = COALESCE(NULLIF($4, ''), producer),
            page_count = CASE WHEN $5::int <= 0 THEN page_count ELSE $5 END,
            updated_at = NOW()
        WHERE id = $6
        RETURNING id, name, isbn, file_path, mime_type, page_count, producer, subject, pdf_version, created_at, updated_at
    )
    SELECT 
        ub.id,
        ub.name AS book_name, 
        COALESCE(a.name, '') as author_name,
        ub.isbn,
        ub.mime_type,
        ub.page_count,
        ub.producer,
        ub.subject,
        ub.pdf_version,
        ub.created_at,
        ub.updated_at
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

	row := q.db.QueryRowContext(ctx, updateBook,
		arg.Title,
		arg.Isbn,
		arg.Subject,
		arg.Producer,
		arg.PageCount,
		id,
	)

	var ub UserBook
	err := row.Scan(
		&ub.ID,
		&ub.BookName,
		&ub.AuthorName,
		&ub.ISBN,
		&ub.MimeType,
		&ub.PageCount,
		&ub.Producer,
		&ub.Subject,
		&ub.PdfVersion,
		&ub.CreatedAt,
		&ub.UpdatedAt,
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
