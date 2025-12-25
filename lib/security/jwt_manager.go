// internal/security/jwt.go
package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret string
	ttl    int64
}

func NewJWTManager(secret string, ttl int64) *JWTManager {
	return &JWTManager{secret: secret, ttl: ttl}
}

func (j *JWTManager) Generate(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Second * time.Duration(j.ttl)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}


