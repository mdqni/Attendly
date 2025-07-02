package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

type CustomClaims struct {
	UserID string   `json:"user-id"`
	Perms  []string `json:"perms"`
	jwt.RegisteredClaims
}

func ParseJWT(tokenStr, secret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		log.Println(err)
		return nil, err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}

func GenerateJWT(secret, userID string, perms []string, ttl time.Duration) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		Perms:  perms,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
