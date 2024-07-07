package model

import "github.com/golang-jwt/jwt"

type Token struct {
	UID     int
	Token   string
	Expires int64
}

type TokenClaims struct {
	UserID     int   `json:"user_id"`
	Authorized bool  `json:"authorized"`
	Expires    int64 `json:"expires"`
	jwt.StandardClaims
}
