package repository

import (
	"context"
	"errors"
	model "loadsg/lib/model"

)

var ErrJWTRefreshNotFound = errors.New("jwt refresh token not found")

type JWTRefreshRepository interface {
	GetByJWTHash(ctx context.Context, hash string) (*model.JWTRefreshToken, error)
	Create(ctx context.Context, jwt *model.JWTRefreshToken) (*model.JWTRefreshToken, error)
}


