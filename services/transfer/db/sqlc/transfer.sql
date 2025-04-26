-- name: CreateTransfer :one
INSERT INTO transfers (
    id,
    from_account_id,
    to_account_id,
    amount
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: SearchTransfersWithAccount :many
SELECT * FROM transfers
WHERE from_account_id = @search_account OR to_account_id = @search_account
ORDER BY created_at DESC
OFFSET $1 LIMIT $2;
