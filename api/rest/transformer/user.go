package transformer

import (
	"github.com/muzz/api/rest/definition"
	"github.com/muzz/api/service/entity"
)

func FromUserInputDefToEntity(in definition.UserInput) entity.UserInput {
	return entity.UserInput{
		Email:    in.Email,
		Password: in.Password,
		Name:     in.Name,
		Gender:   in.Gender,
		DOB:      in.DOB,
	}
}

func FromUserEntityToDef(in entity.User) definition.User {
	return definition.User{
		ID:       in.ID,
		Email:    in.Email,
		Password: in.Password,
		Name:     in.Name,
		Gender:   in.Gender,
		Age:      getAge(in.DOB),
	}
}

func FromTokenEntityToDef(in entity.Token) definition.Token {
	return definition.Token{
		Token:   in.Token,
		Expires: in.Expires,
	}
}
