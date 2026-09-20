-- +goose Up
ALTER TABLE book_authors ALTER COLUMN api_key DROP DEFAULT;

-- +goose Down
ALTER TABLE book_authors ALTER COLUMN api_key SET DEFAULT (
    encode(sha256(random()::text::bytea), 'hex')
);