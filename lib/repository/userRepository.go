package repository

import (
	"context"
	"errors"
	model "loadsg/lib/model"

)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	GetByUID(ctx context.Context, uid string) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
}


