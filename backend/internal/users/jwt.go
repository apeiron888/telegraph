package users

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret []byte
	exp    time.Duration
}

func NewJWTManager(secret string, exp time.Duration) *JWTManager {
	return &JWTManager{
		secret: []byte(secret),
		exp:    exp,
	}
}

func (j *JWTManager) Generate(userID, role string) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(j.exp).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTManager) Verify(tokenStr string) (string, error) {
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return j.secret, nil
	})
	if err != nil || !tok.Valid {
		return "", errors.New("invalid_token")
	}

	claims := tok.Claims.(jwt.MapClaims)
	return claims["user_id"].(string), nil
}
