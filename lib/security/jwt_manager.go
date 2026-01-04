// internal/security/jwt.go
package security

import (
	"time"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret []byte
	refreshSecret []byte
	issuer string
	expires_at int64
}

type Claims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}


func NewJWTManager(secret, refreshSecret, issuer string, expires_at int64) *JWTManager {
	return &JWTManager{secret: []byte(secret), refreshSecret: []byte(refreshSecret), issuer: issuer, expires_at: expires_at}
}


func (j *JWTManager) Generate(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject: userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expires_at))),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Issuer: j.issuer,
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}


func (j *JWTManager) GenerateRefresh(userID string) (string, error){
	claims := &jwt.RegisteredClaims{
		Subject: userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7*24*time.Hour)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Issuer: j.issuer,
	}
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.refreshSecret))
}


func (j *JWTManager) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}
			return j.secret, nil
		},
		jwt.WithIssuer(j.issuer),
		jwt.WithLeeway(2*time.Minute),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

