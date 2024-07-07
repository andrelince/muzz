package transformer

import (
	"github.com/muzz/api/repository/model"
	"github.com/muzz/api/service/entity"
)

func FromUserEntityInputToModel(in entity.UserInput) model.UserInput {
	return model.UserInput{
		Email:    in.Email,
		Password: in.Password,
		Name:     in.Name,
		Gender:   in.Gender,
		DOB:      in.DOB,
	}
}

func FromUserModelToEntity(in model.User) entity.User {
	return entity.User{
		ID:       in.ID,
		Email:    in.Email,
		Password: in.Password,
		Name:     in.Name,
		Gender:   in.Gender,
		DOB:      in.DOB,
	}
}
