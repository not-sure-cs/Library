-- +goose Up
ALTER TABLE books ADD file_path TEXT UNIQUE NOT NULL; 

-- +goose Down
ALTER TABLE books DROP file_path;