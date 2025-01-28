package repository

import (
	"context"
	"time"

	"github.com/myacey/jxgercorp-banking/token/internal/models"
)

type TokenRepository interface {
	CreateToken(c context.Context, payload string, username string, ttl time.Duration) error
	GetToken(c context.Context, username string) (*models.Payload, string, error)
	UpdateToken(c context.Context, newPayload *models.Payload, username string, ttl time.Duration) error
}
