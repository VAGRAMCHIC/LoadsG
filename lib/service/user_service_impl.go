package service

import (
	"context"

	"loadsg/lib/dto"
	"loadsg/lib/model"
	"loadsg/lib/repository"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) Create(ctx context.Context, req dto.CreateUserRequest) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Token), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Comment: req.Comment,
		Token:   string(hash),
	}

	return s.repo.Create(ctx, user)
}
