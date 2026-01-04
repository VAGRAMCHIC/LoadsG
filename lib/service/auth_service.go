package service

import (
	"context"
)

type AuthService interface {
    Login(ctx context.Context, uid string, password string) (token string, refreshToken string, err error)
}

