// internal/service/auth_service_impl.go
package service

import (
	"context"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"

	"loadsg/lib/model"
	"loadsg/lib/repository"
	"loadsg/lib/security"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUnauthorized = errors.New("Unauthorized")
var ErrForbidden = errors.New("Forbiden")

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
	uid string,
	token string,
) (*security.TokenPair, error) {

	// 1. Получаем пользователя
	user, err := s.users.GetByUID(ctx, uid)
	if err != nil {
		log.Printf("user not found: sent_uid=%s", uid)
		return nil, ErrInvalidCredentials
	}

	// 2. Проверяем секрет (пароль / токен)
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Token),
		[]byte(token),
	); err != nil {
		log.Printf("invalid token for uid=%s", uid)
		return nil, ErrInvalidCredentials
	}

	// 3. Генерируем пару токенов
	pair, err := s.jwt.GeneratePair(user.UID)
	if err != nil {
		return nil, err
	}

	// 4. Сохраняем refresh token (ТОЛЬКО ХЕШ)
	_, err = s.tokens.Create(ctx, &model.JWTRefreshToken{
		UserUid:   user.UID,
		TokenHash: security.HashToken(pair.RefreshToken),
		ExpiresAt: pair.RefreshExp,
	})
	if err != nil {
		return nil, err
	}

	// 5. Возвращаем TokenPair
	return pair, nil
}



func (s *authService) Refresh(
	ctx context.Context,
	refreshToken string,
) (*security.TokenPair, error) {

	claims, err := s.jwt.ValidateRefresh(refreshToken)
	if err != nil {
		return nil, ErrUnauthorized
	}

	userID := claims.Subject

	hash := security.HashToken(refreshToken)

	dbUserID, err := s.tokens.Get(ctx, hash)
	if err != nil {
		return nil, ErrForbidden
	}

	// защита: user_id из JWT и из БД должны совпадать
	if dbUserID != userID {
		return nil, ErrForbidden
	}

	_ = s.tokens.Delete(ctx, hash)

	pair, err := s.jwt.GeneratePair(userID)
	if err != nil {
		return nil, err
	}

	err = s.tokens.Save(
		ctx,
		userID,
		security.HashToken(pair.RefreshToken),
		pair.RefreshExp,
	)
	if err != nil {
		return nil, err
	}

	return pair, nil
}


