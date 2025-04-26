package accountrepo

import (
	"github.com/google/uuid"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/models/entity"
)

type SearchAccountsParams struct {
	OwnerUsername string
	Currency      entity.Currency
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
