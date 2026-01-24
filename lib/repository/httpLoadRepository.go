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
	DeleteFixed(ctx context.Context, id string) error

	GetConstantById(ctx context.Context, id string) (*model.ConstantHttpLoad, error)
	CreateConstant(ctx context.Context, httpLoad *model.ConstantHttpLoad) (*model.ConstantHttpLoad, error)
	DeleteConstant(ctx context.Context, id string) error

	GetRampUpById(ctx context.Context, id string) (*model.RampUpHttpLoad, error)
	CreateRampUp(ctx context.Context, httpLoad *model.RampUpHttpLoad) (*model.RampUpHttpLoad, error)

	GetFakeById(ctx context.Context, id string) (*model.FakeHttpLoad, error)
	CreateFake(ctx context.Context, httpLoad *model.FakeHttpLoad) (*model.FakeHttpLoad, error)

	DeleteRampUp(ctx context.Context, id string) error
}
