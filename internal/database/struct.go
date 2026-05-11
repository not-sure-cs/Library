package database

import "database/sql"

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
	BookName   string         `json:"book"`
	AuthorName string         `json:"author"`
	ISBN       sql.NullString `json:"isbn"`
}
