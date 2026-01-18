-- name: CreateTransfer :one
INSERT INTO
    transfers (
        id,
        from_account_id,
        to_account_id,
        currency_code,
        amount,
        from_account_username,
        to_account_username
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: SearchTransfers :many
SELECT *
FROM transfers
WHERE
    -- проверка текущего пользователя
    (
        from_account_id = @current_account_id
        OR to_account_id = @current_account_id
    )

    -- фильтрация по второму участнику (username)
    AND (
        sqlc.narg(with_username)::TEXT IS NULL
        OR from_account_username =sqlc.narg(with_username)
        OR to_account_username = sqlc.narg(with_username)
    )

    -- фильтрация по второму участнику (account_id)
    AND (
        sqlc.narg('with_account_id')::UUID IS NULL
        OR from_account_id = sqlc.narg('with_account_id')
        OR to_account_id = sqlc.narg('with_account_id')
    )

    -- по валютам
    AND (
        sqlc.narg(currency)::TEXT IS NULL
        OR currency_code = sqlc.narg(currency)
    )
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');