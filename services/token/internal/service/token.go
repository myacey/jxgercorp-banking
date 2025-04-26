package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	JwtClaimsUsername = "username"
	JwtClaimExp       = "exp"

	JwtExpireTime = 24 * time.Hour
)

type TokenServiceConfig struct {
	SecretKey string `mapstructure:"jwt_secret_key"`
}

type Token struct {
	secretKey []byte
}

func NewToken(secretKey string) *Token {
	return &Token{secretKey: []byte(secretKey)}
}

func (t *Token) GenerateToken(ctx context.Context, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		JwtClaimsUsername: username,
		JwtClaimExp:       time.Now().Add(JwtExpireTime).Unix(),
	})
	tokenStr, err := token.SignedString(t.secretKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (t *Token) ValidateToken(ctx context.Context, tokenStr string) (map[string]interface{}, error) {
	log.Printf("validate token: %s", tokenStr)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return t.secretKey, nil
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
