// internal/service/auth_service_impl.go
package service

import (
	"context"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"

	"loadsg/lib/repository"
	"loadsg/lib/model"
	"loadsg/lib/security"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type authService struct {
    users repository.UserRepository
    tokens repository.JWTRefreshRepository
		jwt    security.JWTManager
}

func NewAuthService(
    users 	repository.UserRepository,
		tokens 	repository.JWTRefreshRepository,
    jwt 		security.JWTManager,
) AuthService {
    return &authService{users: users, tokens: tokens, jwt: jwt}
}

func (s *authService) Login(
    ctx context.Context,
    uid string, token string,
) (string, string, error) {

    user, err := s.users.GetByUID(ctx, uid)
    if err != nil {
			log.Printf("user: %s, sent_uid: %s", user.UID, uid)
      return "","", ErrInvalidCredentials
    }

    if err := bcrypt.CompareHashAndPassword(
      []byte(user.Token),
      []byte(token),
    ); err != nil {
			log.Printf("token: %s", token)
      return "","", ErrInvalidCredentials
    }

		acces_token, err := s.jwt.Generate(user.UID)
    if err != nil {
      return "","", err
    }
		refreshToken, err := s.jwt.GenerateRefresh(user.UID)
		if err != nil {
			return "","", err
		}
		
		jwtRefreshToken := &model.JWTRefreshToken{
			UserUid : user.UID,
			TokenHash: refreshToken,	
		}
		s.tokens.Create(ctx, jwtRefreshToken)
	return acces_token, refreshToken, nil
}

