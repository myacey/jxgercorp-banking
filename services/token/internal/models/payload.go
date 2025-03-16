package models

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/myacey/jxgercorp-banking/services/shared/cstmerr"
)

var ErrTokenExpired = errors.New("expired token")

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	ExpireAt  time.Time `json:"expire_at"`
}

func NewPayload(username string, expireTime time.Time) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}
	return &Payload{
		ID:        tokenId,
		Username:  username,
		CreatedAt: time.Now(),
		ExpireAt:  expireTime,
	}, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpireAt) {
		return ErrTokenExpired
	}

	return nil
}
