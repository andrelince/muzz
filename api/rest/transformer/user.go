package transformer

import (
	"github.com/muzz/api/rest/definition"
	"github.com/muzz/api/service/entity"
)

func FromUserInputDefToEntity(in definition.UserInput) entity.UserInput {
	return entity.UserInput{
		Email:        in.Email,
		Password:     in.Password,
		Name:         in.Name,
		Gender:       in.Gender,
		DOB:          in.DOB,
		LocationLat:  in.LocationLat,
		LocationLong: in.LocationLong,
	}
}

func FromUserEntityToDef(in entity.User) definition.User {
	return definition.User{
		ID:           in.ID,
		Email:        in.Email,
		Password:     in.Password,
		Name:         in.Name,
		Gender:       in.Gender,
		Age:          getAge(in.DOB),
		LocationLat:  in.LocationLat,
		LocationLong: in.LocationLong,
	}
}

func FromTokenEntityToDef(in entity.Token) definition.Token {
	return definition.Token{
		Token:   in.Token,
		Expires: in.Expires,
	}
}

func FromMatchEntityToDef(in entity.Match) definition.Match {
	out := definition.Match{
		Matched: in.IsMatch,
	}
	if in.IsMatch {
		out.MatchID = &in.ID
	}
	return out
}

func FromDiscoveryEntityToDef(in entity.Discovery) definition.Discovery {
	return definition.Discovery{
		User:                FromUserEntityToDef(in.User),
		DistanceFromMe:      in.DistanceFromMe,
		AttractivenessScore: in.AttractivenessScore,
	}
}
