package request

import "github.com/google/uuid"

type CreateAccount struct {
	OwnerUsername string `json:"-"`
	Currency      string `json:"currency" binding:"required"`
}

type SearchAccounts struct {
	Username string `json:"username" binding:"required"`
	Currency string `json:"currency" binding:"omitempty"`
}

type CreateTransfer struct {
	FromAccountID uuid.UUID `json:"from_account_id"`
	ToAccountID   uuid.UUID `json:"to_account_id"`
	Amount        int64     `json:"amount"`
}

type SearchTransfersWithAccount struct {
	AccountID uuid.UUID `json:"account_id"`
	Offset    int32     `json:"offset"`
	Limit     int32     `json:"limit"`
}

type AddAccountBalance struct {
	ID         uuid.UUID `json:"account_id"`
	AddBalance int64     `json:"add_balance"`
}
