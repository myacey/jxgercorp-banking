package tokenmaker

import (
	"errors"
	"strconv"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/myacey/jxgercorp-banking/token/internal/models"
	"github.com/o1egl/paseto"
)

var (
	ErrSymmetricKeyTooShot = errors.New("symmetric key too short: " + strconv.Itoa(chacha20poly1305.KeySize))
	ErrKeyNotValid         = errors.New("providen key is not valid")
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPaseto(symmetricKey string) (*PasetoMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, ErrSymmetricKeyTooShot
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (m *PasetoMaker) CreateToken(username string, ttl time.Duration) (string, error) {
	payload, err := models.NewPayload(username, ttl)
	if err != nil {
		return "", err
	}

	return m.paseto.Encrypt(m.symmetricKey, payload, nil)
}

func (m *PasetoMaker) VerifyToken(token string) (*models.Payload, error) {
	payload := &models.Payload{}

	if err := m.paseto.Decrypt(token, m.symmetricKey, &payload, nil); err != nil {
		return nil, ErrKeyNotValid
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}

	return payload, nil
}
