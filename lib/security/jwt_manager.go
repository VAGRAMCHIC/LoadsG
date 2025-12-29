// internal/security/jwt.go
package security

import (
	"time"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
)

type JWTManager struct {
	secret []byte
	issuer string
	expires_at int64
}

type Claims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}


func NewJWTManager(secret, issuer string, expires_at int64) *JWTManager {
	return &JWTManager{secret: []byte(secret), issuer: issuer, expires_at: expires_at}
}

func (j *JWTManager) Generate(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": j.expires_at,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}


func setRefreshCookie(c *gin.Context, refreshToken string) {
	c.SetCookie(
		"refresh_token",        // name
		refreshToken,           // value
		int((30 * 24 * time.Hour).Seconds()), // maxAge
		"/auth/refresh",        // path
		"",                     // domain ("" = текущий)
		true,                   // secure (true в prod)
		true,                   // httpOnly
	)
}

func clearRefreshCookie(c *gin.Context) {
	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/auth/refresh",
		"",
		true,
		true,
	)
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

