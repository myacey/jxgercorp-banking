package entity

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/models/dto/response"
)

type Currency string

const (
	CurrencyRUB Currency = "RUB"
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
)

func (s Currency) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *Currency) Scan(value interface{}) error {
	*s = Currency(string(value.(string)))
	return nil
}

type Account struct {
	ID            uuid.UUID
	OwnerUsername string
	Balance       int64
	Currency      Currency
	CreatedAt     time.Time
}

func (a *Account) ToResponse() *response.Account {
	return &response.Account{
		ID:            a.ID,
		OwnerUsername: a.OwnerUsername,
		Balance:       a.Balance,
		Currency:      string(a.Currency),
		CreatedAt:     a.CreatedAt,
	}
}
