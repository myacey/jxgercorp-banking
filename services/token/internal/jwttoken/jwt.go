package jwttoken

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	JwtClaimID   = "uuid"
	JwtClaimRole = "role"
	JwtClaimExp  = "exp"
)

type TokenServiceConfig struct {
	SecretKey string `mapstructure:"jwt_secret_key"`
}

type Service struct {
	secretKey []byte
}

func New(cfg TokenServiceConfig) *Service {
	return &Service{[]byte(cfg.SecretKey)}
}

func (s *Service) CreateDummyToken(role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		JwtClaimRole: role,
		JwtClaimExp:  time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenStr, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (s *Service) CreateUserToken(id uuid.UUID, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		JwtClaimID:   id.String(),
		JwtClaimRole: role,
		JwtClaimExp:  time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenStr, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (s *Service) VerifyToken(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
