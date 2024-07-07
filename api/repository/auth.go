package repository

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/muzz/api/pkg/redis"
	"github.com/muzz/api/repository/model"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -destination=./mocks/mock_user_connector.go -package=mocks github.com/muzz/api/repository UserConnector
type AuthConnector interface {
	HashPassword(value string) (string, error)
	ValidateHash(hashed, value string) error
	GenerateToken(ctx context.Context, uid int) (model.Token, error)
}

type AuthRepo struct {
	l      *logrus.Logger
	cache  *redis.Redis
	secret string
}

func NewAuthRepo(l *logrus.Logger, cache *redis.Redis, secret string) AuthRepo {
	return AuthRepo{
		l:     l,
		cache: cache,
	}
}

func (a AuthRepo) GenerateToken(ctx context.Context, uid int) (model.Token, error) {
	secretKey := []byte(a.secret)
	expires := time.Now().Add(time.Minute * 30).Unix()

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = uid
	claims["exp"] = expires

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secretKey)
	if err != nil {
		return model.Token{}, err
	}

	return model.Token{
		UID:     uid,
		Token:   signed,
		Expires: expires,
	}, nil
}

func (s AuthRepo) ValidateHash(hashed, value string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(value))
}

func (a AuthRepo) HashPassword(value string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	return string(bytes), err
}
