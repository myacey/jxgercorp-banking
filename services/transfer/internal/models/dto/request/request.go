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

type CreateTransfer struct {
	FromAccountID uuid.UUID `json:"from_account_id"`
	ToAccountID   uuid.UUID `json:"to_account_id"`
	Amount        int64     `json:"amount"`
}

type SearchTransfersWithAccount struct {
	AccountID string `form:"account_id" binding:"required,uuid"`
	Offset    int32  `form:"offset" binding:"omitempty,min=0"`
	Limit     int32  `form:"limit" binding:"omitempty,min=1,max=100"`
}

type AddAccountBalance struct {
	ID         uuid.UUID `json:"account_id"`
	AddBalance int64     `json:"add_balance"`
}
