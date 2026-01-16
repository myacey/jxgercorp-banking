package request

import "github.com/google/uuid"

type CreateAccount struct {
	OwnerUsername string `json:"-"`
	Currency      string `json:"currency" binding:"required"`
}

type SearchAccounts struct {
	Username string `form:"username" binding:"required"`
	Currency string `form:"currency" binding:"omitempty"`
}

type DeleteAccount struct {
	AccountID string `form:"account_id" binding:"required"`
}

type CreateTransfer struct {
	FromAccountID       uuid.UUID `json:"from_account_id" binding:"required"`
	FromAccountUsername string    `json:"from_account_username" binding:"required"`
	ToAccountID         uuid.UUID `json:"to_account_id" binding:"required"`
	ToAccountUsername   string    `json:"to_account_username" binding:"required"`
	Amount              float64   `json:"amount" binding:"required"`
	Currency            string    `json:"currency" binding:"required"`
}

type SearchTransfersWithAccount struct {
	CurrentAccountID string `form:"current_account_id" binding:"required,uuid"`
	WithUsername     string `form:"with_username" binding:"omitempty"`
	WithAccountID    string `form:"with_account_id" binding:"omitempty"`
	CurrencyCode     string `form:"currency" binding:"omitempty"`
	Offset           int32  `form:"offset" binding:"omitempty,min=0"`
	Limit            int32  `form:"limit" binding:"omitempty,min=1,max=100"`
}

type AddAccountBalance struct {
	ID         uuid.UUID `json:"account_id"`
	AddBalance float64   `json:"add_balance"`
}
