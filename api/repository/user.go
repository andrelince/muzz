package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
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
	Discover(ctx context.Context, userID int, age []int, gender string) ([]model.Discovery, error)
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
	query := `INSERT INTO users (email, password, name, gender, date_of_birth, location_lat, location_long) 
              VALUES (:email, :password, :name, :gender, :date_of_birth, :location_lat, :location_long) RETURNING *`

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

func (r UserRepo) Discover(ctx context.Context, userID int, age []int, gender string) ([]model.Discovery, error) {
	var results []model.Discovery

	// Get the current user's location
	var currentUser model.User
	err := r.db.DBX().GetContext(ctx, &currentUser, "SELECT location_lat, location_long FROM users WHERE id = $1", userID)
	if err != nil {
		return nil, err
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.Select(
		`u.id AS "user.id"`,
		`u.name AS "user.name"`,
		`u.email AS "user.email"`,
		`u.gender AS "user.gender"`,
		`u.date_of_birth AS "user.date_of_birth"`,
		`u.location_lat AS "user.location_lat"`,
		`u.location_long AS "user.location_long"`,
		"(SELECT COUNT(*) FROM user_swipes s WHERE s.swiped_user_id = u.id AND s.swipe_status = true) AS attractiveness_score",
		`earth_distance(
			ll_to_earth(COALESCE($1, 0), COALESCE($2, 0)), 
			ll_to_earth(COALESCE(u.location_lat, 0), COALESCE(u.location_long, 0))
		) AS distance_from_me`,
	).
		From("users u").
		LeftJoin("matches m1 ON (u.id = m1.user1_id OR u.id = m1.user2_id) AND (m1.user1_id = $1 OR m1.user2_id = $3)").
		LeftJoin("user_swipes s ON u.id = s.swiped_user_id AND s.user_id = $3").
		Where("u.id != $3").
		Where("m1.user1_id IS NULL").
		Where("s.user_id IS NULL").
		OrderBy("distance_from_me", "attractiveness_score DESC")

	args := []interface{}{currentUser.LocationLat, currentUser.LocationLong, userID}

	if len(age) == 2 {
		query = query.Where(fmt.Sprintf("DATE_PART('year', AGE(u.date_of_birth)) BETWEEN %d AND %d", age[0], age[1]))
	}

	if gender != "" {
		query = query.Where(fmt.Sprintf("u.gender = '%s'", gender))
	}

	sql, _, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.DBX().SelectContext(ctx, &results, sql, args...)
	if err != nil {
		return nil, err
	}

	return results, nil
}
