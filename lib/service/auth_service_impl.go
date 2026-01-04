// internal/service/auth_service_impl.go
package service

import (
	"context"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"loadsg/lib/model"
	"loadsg/lib/repository"
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
) (map[string]string, map[string]string, error) {

    user, err := s.users.GetByUID(ctx, uid)
    if err != nil {
			log.Printf("user: %s, sent_uid: %s", user.UID, uid)
			return map[string]string{"":""},map[string]string{"":""}, ErrInvalidCredentials
    }

    if err := bcrypt.CompareHashAndPassword(
      []byte(user.Token),
      []byte(token),
    ); err != nil {
			log.Printf("token: %s", token)
			return map[string]string{"":""},map[string]string{"":""}, ErrInvalidCredentials
    }

		acces_token, acExp, err := s.jwt.Generate(user.UID)
    if err != nil {
			return map[string]string{"":""},map[string]string{"":""}, err
    }
		refreshToken, reExp, err := s.jwt.GenerateRefresh(user.UID)
		if err != nil {
			return map[string]string{"":""},map[string]string{"":""}, err
		}
		
		jwtRefreshToken := &model.JWTRefreshToken{
			UserUid : user.UID,
			TokenHash: refreshToken,
			ExpiresAt: reExp,
		}
		s.tokens.Create(ctx, jwtRefreshToken)
	return map[string]string{acces_token:acExp.Format(time.RFC3339)}, map[string]string{refreshToken:reExp.Format(time.RFC3339)}, nil
}

