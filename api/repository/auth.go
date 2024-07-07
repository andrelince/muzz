package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
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
	GetTokenClaims(ctx context.Context, token string, out any) error
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

	var claims jwt.MapClaims

	err := mapstructure.Decode(model.TokenClaims{
		UserID:     uid,
		Authorized: true,
		Expires:    expires,
	}, &claims)

	if err != nil {
		return model.Token{}, err
	}

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

func (a AuthRepo) GetTokenClaims(ctx context.Context, tokenStr string, out any) error {
	fmt.Printf("Claims: %+v\n", out)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.secret), nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		err := mapstructure.Decode(claims, &out)
		if err != nil {
			return err
		}
	}
	return nil
}
