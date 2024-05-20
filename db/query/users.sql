-- name: CreateUser :one
INSERT INTO users (name) VALUES ($1)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1
LIMIT 1
FOR NO KEY UPDATE;

-- name: GetUserForUpdate :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateImageCount :one
UPDATE users
SET image_count = image_count + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateSubscribed :one
UPDATE users
SET subscribed = $2
WHERE id = $1
RETURNING *;
