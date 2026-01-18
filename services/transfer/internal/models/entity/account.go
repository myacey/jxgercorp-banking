package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/models/dto/response"
)

type Account struct {
	ID            uuid.UUID
	OwnerUsername string
	Balance       int64
	CurrencyCode  string
	CreatedAt     time.Time
}

func (a *Account) ToResponse() *response.Account {
	return &response.Account{
		ID:            a.ID,
		OwnerUsername: a.OwnerUsername,
		Balance:       float64(a.Balance) / 100,
		Currency:      a.CurrencyCode,
		CreatedAt:     a.CreatedAt,
	}
}

type Currency struct {
	Code      string
	Symbol    string
	Precision int
}

func (c *Currency) ToResponse() *response.Currency {
	return &response.Currency{
		Code:      c.Code,
		Symbol:    c.Symbol,
		Precision: c.Precision,
	}
}
