-- name: CreateBook :one
INSERT INTO books (id, created_at, updated_at, name, isbn)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: CreateAuthor :one
INSERT INTO authors (id,created_at, updated_at, name)
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: GetAuthor :one
SELECT *
FROM authors
WHERE name = $1
LIMIT 1;

-- name: LinkBookAuthor :one
INSERT INTO book_authors (book_id, author_id)
VALUES($1, $2)
RETURNING *;

-- name: DeleteBook :exec
DELETE FROM books WHERE id = $1;

-- name: UnlinkBook :exec
DELETE FROM book_authors WHERE book_id = $1;


