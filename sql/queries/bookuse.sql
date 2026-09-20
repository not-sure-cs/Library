-- name: GetBook :one
SELECT 
    books.id,
    books.name AS book_name,
    COALESCE(authors.name, '') AS author_name,
    books.isbn,
    books.file_path,
    books.mime_type,
    books.page_count,
    books.producer,
    books.subject,
    books.pdf_version,
    books.created_at,
    books.updated_at
FROM books
LEFT JOIN book_authors ON books.id = book_authors.book_id
LEFT JOIN authors ON book_authors.author_id = authors.id
WHERE books.id = $1
LIMIT 1;

-- name: GetAuthorBooks :many
SELECT 
    books.id,
    books.name AS book_name,
    COALESCE(authors.name, '') AS author_name,
    books.isbn,
    books.file_path,
    books.mime_type,
    books.page_count,
    books.producer,
    books.subject,
    books.pdf_version,
    books.created_at,
    books.updated_at
FROM book_authors 
JOIN books ON book_authors.book_id = books.id
JOIN authors ON book_authors.author_id = authors.id
WHERE authors.id = $1;

-- name: GetMetaData :one
SELECT file_path, mime_type FROM books
WHERE id = $1
LIMIT 1;

-- name: CheckApiKeyExists :one
SELECT EXISTS (
  SELECT 1 FROM book_authors WHERE api_key = $1
)AS exists;