-- +goose Up
ALTER TABLE books ADD mime_type TEXT; 

-- +goose Down
ALTER TABLE books DROP mime_type;