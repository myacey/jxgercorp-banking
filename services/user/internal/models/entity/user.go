package entity

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
)

type UserUnhashed struct {
	Username string
	Email    string
	Password string
}

type UserStatus string

const (
	UserStatusPending UserStatus = "pending"
	UserStatusActive  UserStatus = "active"
	UserStatusBanned  UserStatus = "banned"
)

func (s UserStatus) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *UserStatus) Scan(value interface{}) error {
	*s = UserStatus(string(value.([]byte)))
	return nil
}

type User struct {
	ID             uuid.UUID
	Username       string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
	Status         UserStatus
}
