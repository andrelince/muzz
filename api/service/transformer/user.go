package transformer

import (
	"github.com/muzz/api/repository/model"
	"github.com/muzz/api/service/entity"
)

func FromUserEntityInputToModel(in entity.UserInput) model.UserInput {
	return model.UserInput{
		Email:        in.Email,
		Password:     in.Password,
		Name:         in.Name,
		Gender:       in.Gender,
		DOB:          in.DOB,
		LocationLat:  in.LocationLat,
		LocationLong: in.LocationLong,
	}
}

func FromUserModelToEntity(in model.User) entity.User {
	return entity.User{
		ID:           in.ID,
		Email:        in.Email,
		Password:     in.Password,
		Name:         in.Name,
		Gender:       in.Gender,
		DOB:          in.DOB,
		LocationLat:  in.LocationLat,
		LocationLong: in.LocationLong,
	}
}

func FromTokenModelToEntity(in model.Token) entity.Token {
	return entity.Token{
		Token:   in.Token,
		Expires: in.Expires,
	}
}

func FromMatchModelToEntity(in model.Match) entity.Match {
	return entity.Match{
		ID:      in.ID,
		User1ID: in.User1ID,
		User2ID: in.User2ID,
		IsMatch: in.IsMatch,
	}
}

func FromDiscoveryModelToEntity(in model.Discovery) entity.Discovery {
	return entity.Discovery{
		User:                FromUserModelToEntity(in.User),
		DistanceFromMe:      in.DistanceFromMe,
		AttractivenessScore: in.AttractivenessScore,
	}
}
