package service

import (
	"context"
	"loadsg/lib/dto"
	"loadsg/lib/model"

)

type UserService interface {
    Create(ctx context.Context, req dto.CreateUserRequest) (*model.User, error)
}

