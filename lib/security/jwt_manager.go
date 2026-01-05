// internal/security/jwt.go
package security

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	RefreshExp   time.Time
}



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


func (j *JWTManager) GeneratePair(userID string) (*TokenPair, error) {
	accessExp := time.Now().Add(15 * time.Minute)
	refreshExp := time.Now().Add(7 * 24 * time.Hour)

	accessClaims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(accessExp),
		Issuer:    j.issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	refreshClaims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(refreshExp),
		Issuer:    j.issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	access, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).
		SignedString(j.secret)

	refresh, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).
		SignedString(j.refreshSecret)

	return &TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
		RefreshExp:   refreshExp,
	}, nil
}


func (j *JWTManager) Generate(userID string) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(j.expires_at))

	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    j.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(j.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}



func (j *JWTManager) GenerateRefresh(
	userID string,
) (string, time.Time, error) {

	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    j.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(j.refreshSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}


func (j *JWTManager) ValidateRefresh(
	tokenString string,
) (*jwt.RegisteredClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}
			return j.refreshSecret, nil
		},
		jwt.WithIssuer(j.issuer),
		jwt.WithLeeway(2*time.Minute),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	return claims, nil
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

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}



