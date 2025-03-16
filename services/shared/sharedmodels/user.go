package sharedmodels

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}
