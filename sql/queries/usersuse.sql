-- name: GetUser :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: GetPassHash :one
SELECT pass_hash FROM secrets 
JOIN users ON secrets.user_id = users.id
WHERE users.email = $1
LIMIT 1;