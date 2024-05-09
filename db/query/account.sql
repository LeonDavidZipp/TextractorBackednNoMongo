-- name: CreateAccount :one
INSERT INTO accounts (
    owner,
    email,
    google_id,
    facebook_id
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1
LIMIT 1
FOR NO KEY UPDATE;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1
LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY $1
LIMIT $2
OFFSET $3;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;

-- name: AddImages :one
UPDATE accounts
SET image_count = image_count + sqlc.arg(amount)
RETURNING *;

-- name: UpdateAccount :one
UPDATE accounts
SET email = $2
WHERE id = $1
RETURNING *;
