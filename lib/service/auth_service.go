package service

import (
	"context"
)

type AuthService interface {
    Login(ctx context.Context, uid string, password string) (token map[string]string, refreshToken map[string]string, err error)
}

