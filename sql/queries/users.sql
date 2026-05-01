-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, first_name, last_name, ph_no, email)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: LinkHash :exec
INSERT INTO secrets (user_id, pass_hash)
VALUES($1, $2);

