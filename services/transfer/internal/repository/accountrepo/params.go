package accountrepo

import (
	"github.com/google/uuid"
)

type SearchAccountsParams struct {
	OwnerUsername string
	Currency      string
}

type AddAccountBalance struct {
	AccountID  uuid.UUID
	AddBalance int64
}

type AddTwoAccountsBalance struct {
	FromAccountID uuid.UUID
	ToAccountUD   uuid.UUID
	Amount        int64
}
