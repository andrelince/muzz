package service

import "github.com/muzz/api/repository"

type UserConnector interface {
}

type UserService struct {
	userRepo repository.UserConnector
}

func NewUserService(userRepo repository.UserConnector) UserService {
	return UserService{
		userRepo: userRepo,
	}
}
