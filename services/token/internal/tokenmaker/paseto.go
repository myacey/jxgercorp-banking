package tokenmaker

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/myacey/jxgercorp-banking/services/token/internal/models"
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

func (m *PasetoMaker) CreateToken(username string, expireTime time.Time) (string, error) {
	payload, err := models.NewPayload(username, expireTime)
	if err != nil {
		return "", err
	}

	footer := username // public username

	return m.paseto.Encrypt(m.symmetricKey, payload, footer)
}

// VerifyToken extracts token's payload and footer(public username):
func (m *PasetoMaker) VerifyToken(token string) (*models.Payload, string, error) {
	payload := &models.Payload{}
	var footer string // username

	if err := m.paseto.Decrypt(token, m.symmetricKey, &payload, &footer); err != nil {
		log.Print(err)
		return nil, "", ErrKeyNotValid
	}

	if err := payload.Valid(); err != nil {
		return nil, "", err
	}

	return payload, footer, nil
}
