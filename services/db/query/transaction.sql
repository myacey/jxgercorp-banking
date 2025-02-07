-- name: CreateTransaction :one
INSERT INTO transactions (
    from_user,
    to_user,
    amount
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: SearchOutcomeTransactions :many
SELECT * FROM transactions
WHERE from_user = $1
OFFSET $2 LIMIT $3;

-- name: SearchIncomeTransactions :many
SELECT * FROM transactions
WHERE to_user = $1
OFFSET $2 LIMIT $3;

-- name: SearchTransactionsWithUser :many
SELECT * FROM transactions
WHERE to_user = @search_user OR from_user = @search_user
ORDER BY created_at DESC
OFFSET $1 LIMIT $2;
