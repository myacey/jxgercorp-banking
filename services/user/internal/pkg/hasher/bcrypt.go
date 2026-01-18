package hasher

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	// maxConcurrentBcryptOps represents a max bcrypt operation count.
	maxConcurrentBcryptOps = 20
	// bcryptCost represents cost of bcrypt hashing
	bcryptCost = bcrypt.MinCost
)

// if create new hashers, than create new .go file with all these errors
var (
	ErrTooShort    = errors.New("msg too short")
	ErrTooLong     = errors.New("msg too long")
	ErrDontCompare = errors.New("hash and msg dont match")
)

// Bcrypt realizes hasher with bcrypt.
type Bcrypt struct {
	semaphore chan struct{}
}

func NewBcrypt() *Bcrypt {
	return &Bcrypt{
		semaphore: make(chan struct{}, maxConcurrentBcryptOps),
	}
}

func (h *Bcrypt) Generate(pswd string) (string, error) {
	h.semaphore <- struct{}{}
	defer func() { <-h.semaphore }()

	bytes, err := bcrypt.GenerateFromPassword([]byte(pswd), bcryptCost)
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrHashTooShort):
			return "", ErrTooShort
		case errors.Is(err, bcrypt.ErrPasswordTooLong):
			return "", ErrTooLong
		default:
			return "", err
		}
	}

	return string(bytes), nil
}

func (h *Bcrypt) Compare(hash, msg string) error {
	h.semaphore <- struct{}{}
	defer func() { <-h.semaphore }()

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(msg))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return ErrDontCompare
		default:
			return err
		}
	}

	return nil
}
