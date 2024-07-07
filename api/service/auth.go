package service

import (
	"context"

	"github.com/muzz/api/repository"
)

type AuthConnector interface {
	GetTokenClaims(ctx context.Context, token string, out any) error
}

type AuthService struct {
	authRepo repository.AuthConnector
}

func NewAuthService(authRepo repository.AuthConnector) AuthService {
	return AuthService{
		authRepo: authRepo,
	}
}

func (s AuthService) GetTokenClaims(ctx context.Context, token string, out any) error {
	return s.authRepo.GetTokenClaims(ctx, token, out)
}
