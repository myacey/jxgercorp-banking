-- name: CreateAccount :one
INSERT INTO accounts (
    id,
    owner_username,
    balance,
    currency_code,
    created_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: SearchAccounts :many
SELECT *
FROM accounts
WHERE owner_username = sqlc.arg(username)
  AND (sqlc.narg(currency_code)::text IS NULL OR currency_code = sqlc.narg(currency_code)::text);

-- name: DeleteAccount :one
DELETE FROM accounts
WHERE id  = $1
RETURNING id;

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

-- name: GetCurrencies :many
SELECT * FROM currencies;
