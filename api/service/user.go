package service

import (
	"context"

	"github.com/muzz/api/repository"
	"github.com/muzz/api/service/entity"
	"github.com/muzz/api/service/transformer"
)

type UserConnector interface {
	CreateUser(ctx context.Context, user entity.UserInput) (entity.User, error)
}

type UserService struct {
	userRepo repository.UserConnector
}

func NewUserService(userRepo repository.UserConnector) UserService {
	return UserService{
		userRepo: userRepo,
	}
}

func (s UserService) CreateUser(ctx context.Context, user entity.UserInput) (entity.User, error) {
	in := transformer.FromUserEntityInputToModel(user)

	userM, err := s.userRepo.CreateUser(ctx, in)
	if err != nil {
		return entity.User{}, err
	}

	return transformer.FromUserModelToEntity(userM), nil
}
