package repository

import (
	"context"
	"errors"
	model "loadsg/lib/model"
)

var ErrLoadJobNotFound = errors.New("Load job not found")

type LoadRepository interface {
	GetById(ctx context.Context, id string) (*model.LoadJob, error)
	Create(ctx context.Context, loadJob *model.LoadJob) (*model.LoadJob, error)
	Delete(ctx context.Context, id string) error
}
