package service

import (
	"context"

	"github.com/muzz/api/repository"
	"github.com/muzz/api/service/entity"
	"github.com/muzz/api/service/transformer"
)

type UserConnector interface {
	CreateUser(ctx context.Context, user entity.UserInput) (entity.User, error)
	Login(ctx context.Context, email, password string) (entity.Token, error)
}

type UserService struct {
	userRepo repository.UserConnector
	authRepo repository.AuthConnector
}

func NewUserService(userRepo repository.UserConnector, authRepo repository.AuthConnector) UserService {
	return UserService{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

func (s UserService) CreateUser(ctx context.Context, user entity.UserInput) (entity.User, error) {
	in := transformer.FromUserEntityInputToModel(user)

	// hash password before storing
	hashed, err := s.authRepo.HashPassword(in.Password)
	if err != nil {
		return entity.User{}, err
	}

	in.Password = hashed

	userM, err := s.userRepo.CreateUser(ctx, in)
	if err != nil {
		return entity.User{}, err
	}

	return transformer.FromUserModelToEntity(userM), nil
}

func (s UserService) Login(ctx context.Context, email, password string) (entity.Token, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return entity.Token{}, err
	}

	hashed, err := s.authRepo.HashPassword(password)
	if err != nil {
		return entity.Token{}, err
	}

	if err := s.authRepo.ValidateHash(hashed, password); err != nil {
		return entity.Token{}, err
	}

	token, err := s.authRepo.GenerateToken(ctx, int(user.ID))
	if err != nil {
		return entity.Token{}, err
	}

	return transformer.FromTokenModelToEntity(token), nil
}
