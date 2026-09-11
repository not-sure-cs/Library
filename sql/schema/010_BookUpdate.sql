-- +goose up
ALTER TABLE books DROP genre;
ALTER TABLE books DROP category;
ALTER TABLE books ADD category_code TEXT NOT NULL DEFAULT '0';
-- +goose down
ALTER TABLE books ADD genre TEXT NOT NULL;
ALTER TABLE books ADD category TEXT NOT NULL;
ALTER TABLE books DROP category_code;
