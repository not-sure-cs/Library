package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Linked struct {
	Author   interface{} `json:"author"`
	Book     interface{} `json:"book"`
	Link     interface{} `json:"link"`
	Metadata interface{} `json:"metadata"`
}

type Stream struct {
	Book interface{} `json:"book"`
	Link string      `json:"link"`
}

type UserBook struct {
	ID         uuid.UUID      `json:"id"`
	BookName   string         `json:"book"`
	AuthorName string         `json:"author"`
	ISBN       sql.NullString `json:"isbn"`
	MimeType   sql.NullString `json:"mime_type"`
	PageCount  sql.NullInt32  `json:"page_count"`
	Producer   sql.NullString `json:"producer"`
	Subject    sql.NullString `json:"subject"`
	PdfVersion sql.NullString `json:"pdf_version"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}
