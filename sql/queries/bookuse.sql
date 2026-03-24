-- name: GetBook :many
SELECT books.name,authors.name,api_key FROM book_authors 
JOIN books ON book_authors.book_id = books.id
JOIN authors ON book_authors.author_id = author_id
WHERE books.name = $1
LIMIT 1;
-- name: GetAuthorBook :many
SELECT books.name,authors.name,api_key FROM book_authors 
JOIN books ON book_authors.book_id = books.id
JOIN authors ON book_authors.author_id = author_id
WHERE authors.name = $1;

