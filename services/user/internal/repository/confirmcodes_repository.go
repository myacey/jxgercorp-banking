package repository

import "context"

type ConfirmCodesRepository interface {
	CreateCode(ctx context.Context, username, code string) error
	GetCode(ctx context.Context, username string) (string, error)
}
