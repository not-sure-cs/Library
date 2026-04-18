-- name: GetBook :one
SELECT books.name,authors.name,isbn,books.created_at,books.updated_at,book_id FROM book_authors 
JOIN books ON book_authors.book_id = books.id
JOIN authors ON book_authors.author_id = author_id
WHERE books.id = $1
LIMIT 1;
-- name: GetAuthorBooks :many
SELECT books.name,authors.name,isbn,books.created_at,books.updated_at,book_id,author_id FROM book_authors 
JOIN books ON book_authors.book_id = books.id
JOIN authors ON book_authors.author_id = author_id
WHERE authors.id = $1;

