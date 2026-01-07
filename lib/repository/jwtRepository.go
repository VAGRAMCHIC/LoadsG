package repository

import (
	"context"
	"errors"
	model "loadsg/lib/model"
	"time"
)

var ErrJWTRefreshNotFound = errors.New("jwt refresh token not found")

type JWTRefreshRepository interface {
	GetByJWTHash(ctx context.Context, hash string) (*model.JWTRefreshToken, error)
	Create(ctx context.Context, jwt *model.JWTRefreshToken) (*model.JWTRefreshToken, error)
	Delete(ctx context.Context, jwtHash string) error
	Get(ctx context.Context, hash string) (string, error) // userID
	Save(ctx context.Context, userID string, hash string, expiresAt time.Time) error
}
