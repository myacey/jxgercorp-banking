package sharedmodels

import "time"

type Entry struct {
	WithUser  string    `json:"with_user"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
