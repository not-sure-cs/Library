-- +goose Up
ALTER TABLE books ADD COLUMN cover_path TEXT;
UPDATE books SET cover_path = 'Assets/Covers/Cover-' || id || '.jpg' WHERE cover_path IS NULL;
ALTER TABLE books ALTER COLUMN cover_path SET NOT NULL;
ALTER TABLE books ADD CONSTRAINT books_cover_path_key UNIQUE (cover_path);

-- +goose Down
ALTER TABLE books DROP COLUMN cover_path;