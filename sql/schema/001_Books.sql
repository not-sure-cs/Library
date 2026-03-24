-- +goose Up
CREATE TABLE books (
    id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    isbn TEXT UNIQUE,
    PRIMARY KEY(id)
);

CREATE TABLE authors (
    id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE book_authors (
    book_id UUID REFERENCES books(id) ON DELETE CASCADE,
    author_id UUID REFERENCES authors(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, author_id)
);

-- +goose Down
DROP TABLE book_authors;
DROP TABLE books;
DROP TABLE authors;


