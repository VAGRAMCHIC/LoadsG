package service

import (
	"context"
	"loadsg/lib/security"
)

type AuthService interface {
    Login(ctx context.Context, uid string, password string) (pair *security.TokenPair, err error)
		Refresh(ctx context.Context, refreshToken string) (pair *security.TokenPair, err error)
}

