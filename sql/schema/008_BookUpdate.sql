-- +goose Up
ALTER TABLE books ADD genre TEXT NOT NULL;
ALTER TABLE books ADD category TEXT NOT NULL;
ALTER TABLE books ADD pub_year  SMALLINT NOT NULL; 

-- +goose Down
ALTER TABLE books DROP genre, category, pub_year;