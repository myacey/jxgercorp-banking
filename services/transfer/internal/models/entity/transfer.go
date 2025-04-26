package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/models/dto/response"
)

type Transfer struct {
	ID            uuid.UUID
	FromAccountID uuid.UUID
	ToAccountID   uuid.UUID
	Amount        int64
	CreatedAt     time.Time
}

func (t *Transfer) ToResponse() *response.Transfer {
	return &response.Transfer{
		ID:            t.ID,
		FromAccountID: t.FromAccountID,
		ToAccountID:   t.ToAccountID,
		Amount:        t.Amount,
		CreatedAt:     t.CreatedAt,
	}
}
