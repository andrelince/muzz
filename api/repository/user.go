package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/muzz/api/pkg/pg"
	"github.com/muzz/api/repository/model"
	"github.com/sirupsen/logrus"
)

const (
	uniqueConstraintCode pq.ErrorCode = "23505"
)

var (
	ErrSwipeAlreadyExists = errors.New("swipe already exists")
)

//go:generate mockgen -destination=./mocks/mock_user_connector.go -package=mocks github.com/muzz/api/repository UserConnector
type UserConnector interface {
	CreateUser(ctx context.Context, user model.UserInput) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	Swipe(ctx context.Context, userID, swipedUserID int, status bool) (model.Match, error)
	Discover(ctx context.Context, userID int) ([]model.User, error)
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

func (r UserRepo) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var out model.User

	query := `SELECT * FROM users WHERE email = $1`
	if err := r.db.DBX().Get(&out, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, err
	}
	return out, nil
}

func (r UserRepo) Swipe(ctx context.Context, userID, swipedUserID int, status bool) (model.Match, error) {
	tx, err := r.db.DBX().Beginx()
	if err != nil {
		return model.Match{}, err
	}

	defer func(tx *sqlx.Tx) {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}(tx)

	swipe := model.Swipe{
		UserID:       userID,
		SwipedUserID: swipedUserID,
		SwipeStatus:  status,
		CreatedAt:    time.Now(),
	}

	_, err = tx.NamedExecContext(ctx, `INSERT INTO user_swipes (user_id, swiped_user_id, swipe_status, created_at) 
                                VALUES (:user_id, :swiped_user_id, :swipe_status, :created_at)
                                ON CONFLICT (user_id, swiped_user_id) DO UPDATE SET swipe_status = :swipe_status`, map[string]interface{}{
		"user_id":        swipe.UserID,
		"swiped_user_id": swipe.SwipedUserID,
		"swipe_status":   swipe.SwipeStatus,
		"created_at":     swipe.CreatedAt,
	})
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == uniqueConstraintCode {
			return model.Match{}, ErrSwipeAlreadyExists
		}
		return model.Match{}, err
	}

	var count int
	err = tx.GetContext(ctx, &count, `SELECT COUNT(*) FROM user_swipes 
                                      WHERE user_id IN ($1, $2) AND swiped_user_id IN ($1, $2) AND swipe_status = true`, userID, swipedUserID)
	if err != nil {
		return model.Match{}, err
	}

	if count == 2 {
		match := model.Match{
			User1ID:   userID,
			User2ID:   swipedUserID,
			CreatedAt: time.Now(),
		}
		err = tx.GetContext(ctx, &match, `INSERT INTO matches (user1_id, user2_id, created_at) 
                                           VALUES ($1, $2, $3) 
                                           ON CONFLICT (user1_id, user2_id) DO NOTHING 
                                           RETURNING id, user1_id, user2_id, created_at`, match.User1ID, match.User2ID, match.CreatedAt)
		if err != nil {
			fmt.Println("syntax err", err)
			return model.Match{}, err
		}

		match.IsMatch = true
		return match, nil
	}

	return model.Match{}, nil
}

func (r UserRepo) Discover(ctx context.Context, userID int) ([]model.User, error) {
	var users []model.User

	query := `
		SELECT u.* 
		FROM users u
		LEFT JOIN matches m1 ON (u.id = m1.user1_id OR u.id = m1.user2_id) AND (m1.user1_id = $1 OR m1.user2_id = $1)
		LEFT JOIN user_swipes s ON u.id = s.swiped_user_id AND s.user_id = $1
		WHERE u.id != $1 
		AND m1.user1_id IS NULL 
		AND s.user_id IS NULL
	`

	err := r.db.DBX().SelectContext(ctx, &users, query, userID)
	if err != nil {
		return nil, err
	}

	return users, nil
}
