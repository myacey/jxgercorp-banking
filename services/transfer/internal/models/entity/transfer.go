package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/models/dto/response"
)

type Transfer struct {
	ID uuid.UUID

	FromAccountID       uuid.UUID
	FromAccountUsername string

	ToAccountID       uuid.UUID
	ToAccountUsername string

	Amount       int64
	CurrencyCode string
	CreatedAt    time.Time
}

func (t *Transfer) ToResponse() *response.Transfer {
	return &response.Transfer{
		ID:                  t.ID,
		FromAccountID:       t.FromAccountID,
		FromAccountUsername: t.FromAccountUsername,
		ToAccountUsername:   t.ToAccountUsername,
		ToAccountID:         t.ToAccountID,
		Amount:              float64(t.Amount) / 100,
		Currency:            t.CurrencyCode,
		CreatedAt:           t.CreatedAt,
	}
}
