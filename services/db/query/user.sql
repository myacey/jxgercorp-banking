-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    hashed_password
) VALUES (
    $1, $2, $3
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
    email=COALESCE(sqlc.narg(email), email),
    pending=COALESCE(sqlc.narg(pending), pending)
WHERE
    username=$1
RETURNING *;

-- name: DeleteUserByUsername :exec
DELETE FROM users
WHERE username = $1;

-- name: UpdateTwoUserBalance :many
UPDATE users
SET balance = CASE
    WHEN username = sqlc.arg(from_username) THEN balance - $1
    WHEN username = sqlc.arg(to_username) THEN balance + $1
END
WHERE username IN (sqlc.arg(from_username), sqlc.arg(to_username))
RETURNING username, balance;
