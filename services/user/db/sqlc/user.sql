-- name: CreateUser :one
INSERT INTO users (
    id,
    username,
    email,
    hashed_password
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: UpdateUserInfo :one
UPDATE users
SET
    hashed_password=COALESCE(sqlc.narg(hashed_password), hashed_password),
    email=COALESCE(sqlc.narg(email), email)
WHERE
    username=$1
RETURNING *;

-- name: UpdateUserStatus :one
UPDATE users
SET status=$2
WHERE username=$1
RETURNING *;

-- name: DeleteUserByUsername :exec
DELETE FROM users
WHERE username = $1;
