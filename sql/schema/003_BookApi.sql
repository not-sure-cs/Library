-- +goose Up
ALTER TABLE book_authors ADD column api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT (
    encode(sha256(random()::text::bytea), 'hex')
);
-- +goose Down
ALTER TABLE book_authors DROP column api_key;