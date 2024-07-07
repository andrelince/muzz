package repository

import "github.com/muzz/api/pkg/pg"

//go:generate mockgen -destination=./mocks/mock_user_connector.go -package=mocks github.com/muzz/api/repository UserConnector
type UserConnector interface {
}

type UserRepo struct {
	db *pg.Postgres
}

func NewUserRepo(db *pg.Postgres) UserRepo {
	return UserRepo{
		db: db,
	}
}
