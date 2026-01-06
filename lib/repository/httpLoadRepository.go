package repository

import (
	"context"
	"errors"
	model "loadsg/lib/model"
)

var ErrHTTPLoadNotFound = errors.New("Http load profile not found")

type HttpLoadRepository interface {
	GetFixedById(ctx context.Context, id string) (*model.FixedHttpLoad, error)
	CreateFixed(ctx context.Context, httpLoad *model.FixedHttpLoad) (*model.FixedHttpLoad, error)
}
