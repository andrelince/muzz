package repository

import (
	"context"

	"github.com/muzz/api/pkg/pg"
	"github.com/muzz/api/repository/model"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -destination=./mocks/mock_user_connector.go -package=mocks github.com/muzz/api/repository UserConnector
type UserConnector interface {
	CreateUser(ctx context.Context, user model.UserInput) (model.User, error)
}

type UserRepo struct {
	l  *logrus.Logger
	db *pg.Postgres
}

func NewUserRepo(l *logrus.Logger, db *pg.Postgres) UserRepo {
	return UserRepo{
		l:  l,
		db: db,
	}
}

func (r UserRepo) CreateUser(ctx context.Context, in model.UserInput) (model.User, error) {
	r.l.Info("creating user")
	query := `INSERT INTO users (email, password, name, gender, date_of_birth) 
              VALUES (:email, :password, :name, :gender, :date_of_birth) RETURNING *`

	var out model.User
	stmt, err := r.db.DBX().PrepareNamed(query)
	if err != nil {
		return out, err
	}
	defer stmt.Close()

	err = stmt.Get(&out, in)
	if err != nil {
		return out, err
	}

	return out, nil
}
