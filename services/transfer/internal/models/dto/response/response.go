package response

import (
	"time"

	"github.com/google/uuid"
)

type Error struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Code      int    `json:"code"`
}

type Account struct {
	ID            uuid.UUID `json:"id"`
	OwnerUsername string    `json:"owner_username"`
	Balance       float64   `json:"balance"`
	Currency      string    `json:"currency"`
	CreatedAt     time.Time `json:"created_at"`
}

type Transfer struct {
	ID uuid.UUID `json:"id"`

	FromAccountID       uuid.UUID `json:"from_account_id"`
	FromAccountUsername string    `json:"from_account_username"`

	ToAccountID       uuid.UUID `json:"to_account_id"`
	ToAccountUsername string    `json:"to_account_username"`

	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type Currency struct {
	Code      string `json:"code"`
	Symbol    string `json:"symbol"`
	Precision int    `json:"precision"`
}
