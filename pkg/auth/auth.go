package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenService interface {
	GenerateToken(userID uint) (string, error)
	VerifyToken(tokenStr string) (*jwt.Token, error)
}

type JwtManager struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJwtManager(secretKet string, duration time.Duration) TokenService {
	return &JwtManager{
		secretKey:     secretKet,
		tokenDuration: duration,
	}
}

func (j *JwtManager) GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(j.tokenDuration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JwtManager) VerifyToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
}
