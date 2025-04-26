-- name: CreateAccount :one
INSERT INTO accounts (
    id,
    owner_username,
    balance,
    currency,
    created_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: SearchAccounts :many
SELECT *
FROM accounts
WHERE owner_username = sqlc.arg(username)
  AND (sqlc.narg(currency)::currency_enum IS NULL OR currency = sqlc.narg(currency)::currency_enum);

-- name: GetAccountByID :one
SELECT * FROM accounts
WHERE id = $1
LIMIT 1;

-- name: AddTwoAccountsBalance :many
UPDATE accounts
SET balance = CASE
    WHEN id = @from_account_id THEN balance - $1
    WHEN id = @to_account_id THEN balance + $1
END
WHERE id IN (@from_account_id, @to_account_id)
RETURNING *;


-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + @add_balance
WHERE id = $1
RETURNING *;
