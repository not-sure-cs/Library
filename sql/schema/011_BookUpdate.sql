-- +goose Up
ALTER TABLE books DROP COLUMN category_code;
ALTER TABLE books DROP COLUMN pub_year;
ALTER TABLE books ADD COLUMN page_count INT;
ALTER TABLE books ADD COLUMN producer TEXT;
ALTER TABLE books ADD COLUMN subject TEXT;
ALTER TABLE books ADD COLUMN pdf_version TEXT;

-- +goose Down
ALTER TABLE books DROP COLUMN page_count;
ALTER TABLE books DROP COLUMN producer;
ALTER TABLE books DROP COLUMN subject;
ALTER TABLE books DROP COLUMN pdf_version;
ALTER TABLE books ADD COLUMN category_code TEXT NOT NULL DEFAULT '0';
ALTER TABLE books ADD COLUMN pub_year SMALLINT;
