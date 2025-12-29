// internal/service/auth_service_impl.go
package service

import (
	"context"
	"errors"
	"log"

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

    user, err := s.users.GetByUID(ctx, uid)
    if err != nil {
				log.Printf("user: %s, sent_uid: %s", user.UID, uid)
        return "", ErrInvalidCredentials
    }

    if err := bcrypt.CompareHashAndPassword(
        []byte(user.PasswordHash),
        []byte(password),
    ); err != nil {
				log.Printf("pass: %s",password)
        return "", ErrInvalidCredentials
    }

		token, err := s.jwt.Generate(user.UID)
    if err != nil {
        return "", err
    }
		
    return token, nil
}

