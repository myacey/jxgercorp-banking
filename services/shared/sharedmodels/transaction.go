package sharedmodels

import "time"

type Trx struct {
	ID        int64     `json:"id"`
	FromUser  string    `json:"from_user"`
	ToUser    string    `json:"to_user"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
