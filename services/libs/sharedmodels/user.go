package sharedmodels

import (
	"database/sql/driver"
	"time"
)

type User struct {
	ID        int64      `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Balance   int64      `json:"balance"`
	CreatedAt time.Time  `json:"created_at"`
	Status    UserStatus `json:"status"`
}

type UserStatus string

const (
	UserStatusPending UserStatus = "pending"
	UserStatusActive  UserStatus = "active"
	UserStatusBanned  UserStatus = "banned"
)

// Value реализует интерфейс driver.Valuer
func (us UserStatus) Value() (driver.Value, error) {
	return string(us), nil
}

// Scan реализует интерфейс sql.Scanner
func (us *UserStatus) Scan(value interface{}) error {
	*us = UserStatus(string(value.([]byte)))
	return nil
}
