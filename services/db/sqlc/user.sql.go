// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const changeUserBalance = `-- name: ChangeUserBalance :one
UPDATE users
SET
    balance = balance + $2::BIGINT
WHERE username = $1
RETURNING id, username, email, hashed_password, balance, created_at, pending
`

type ChangeUserBalanceParams struct {
	Username   string `json:"username"`
	AddBalance int64  `json:"add_balance"`
}

func (q *Queries) ChangeUserBalance(ctx context.Context, arg ChangeUserBalanceParams) (User, error) {
	row := q.db.QueryRowContext(ctx, changeUserBalance, arg.Username, arg.AddBalance)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.HashedPassword,
		&i.Balance,
		&i.CreatedAt,
		&i.Pending,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    hashed_password
) VALUES (
    $1, $2, $3
) RETURNING id, username, email, hashed_password, balance, created_at, pending
`

type CreateUserParams struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.HashedPassword,
		&i.Balance,
		&i.CreatedAt,
		&i.Pending,
	)
	return i, err
}

const deleteUserByUsername = `-- name: DeleteUserByUsername :exec
DELETE FROM users
WHERE username = $1
`

func (q *Queries) DeleteUserByUsername(ctx context.Context, username string) error {
	_, err := q.db.ExecContext(ctx, deleteUserByUsername, username)
	return err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, username, email, hashed_password, balance, created_at, pending FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.HashedPassword,
		&i.Balance,
		&i.CreatedAt,
		&i.Pending,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, email, hashed_password, balance, created_at, pending FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.HashedPassword,
		&i.Balance,
		&i.CreatedAt,
		&i.Pending,
	)
	return i, err
}

const updateUserInfo = `-- name: UpdateUserInfo :one
UPDATE users
SET
    hashed_password=COALESCE($2, hashed_password),
    email=COALESCE($3, email),
    pending=COALESCE($4, pending)
WHERE
    username=$1
RETURNING id, username, email, hashed_password, balance, created_at, pending
`

type UpdateUserInfoParams struct {
	Username       string         `json:"username"`
	HashedPassword sql.NullString `json:"hashed_password"`
	Email          sql.NullString `json:"email"`
	Pending        sql.NullBool   `json:"pending"`
}

func (q *Queries) UpdateUserInfo(ctx context.Context, arg UpdateUserInfoParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserInfo,
		arg.Username,
		arg.HashedPassword,
		arg.Email,
		arg.Pending,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.HashedPassword,
		&i.Balance,
		&i.CreatedAt,
		&i.Pending,
	)
	return i, err
}
