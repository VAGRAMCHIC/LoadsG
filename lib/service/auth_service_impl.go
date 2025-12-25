// internal/service/auth_service_impl.go
package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"loadsg/lib/repository"

	"loadsg/lib/security"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type authService struct {
    users repository.UserRepository
    jwt    security.JWTManager
}

func NewAuthService(
    users repository.UserRepository,
    jwt security.JWTManager,
) AuthService {
    return &authService{users: users, jwt: jwt}
}

func (s *authService) Login(
    ctx context.Context,
    uid string, password string,
) (string, error) {

    user, err := s.users.GetById(ctx, uid)
    if err != nil {
        return "", ErrInvalidCredentials
    }

    if err := bcrypt.CompareHashAndPassword(
        []byte(user.PasswordHash),
        []byte(password),
    ); err != nil {
        return "", ErrInvalidCredentials
    }

		token, err := s.jwt.Generate(user.UID)
    if err != nil {
        return "", err
    }

    return token, nil
}

