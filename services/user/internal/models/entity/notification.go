package entity

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID              `json:"id"`
	UserID    uuid.UUID              `json:"user_id"`
	Username  string                 `json:"username"`
	Email     string                 `json:"email"`
	Type      string                 `json:"type"`
	Subject   string                 `json:"subject"`
	Text      string                 `json:"text"`
	Payload   map[string]interface{} `json:"payload"`
	CreatedAt time.Time              `json:"created_at"`
}
