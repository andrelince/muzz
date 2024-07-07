package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/muzz/api/service"
)

type contextKey string

const userIDKey = contextKey("userID")

type TokenClaims struct {
	UserID     int   `json:"user_id"`
	Authorized bool  `json:"authorized"`
	Expires    int64 `json:"expires"`
}

type AuthMiddleware interface {
	Handle(next http.Handler) http.Handler
}

type AuthHandler struct {
	authService service.AuthConnector
}

func NewAuthHandler(authService service.AuthConnector) AuthHandler {
	return AuthHandler{
		authService: authService,
	}
}

func (m AuthHandler) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "invalid token format", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		var claims TokenClaims
		if err := m.authService.GetTokenClaims(ctx, tokenString, &claims); err != nil {
			http.Error(w, "failed to retrieve token claims", http.StatusUnauthorized)
			return
		}

		fmt.Println(claims.UserID)

		ctx = context.WithValue(ctx, userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(userIDKey).(int)
	if !ok {
		return 0, errors.New("user id not found")
	}
	return userID, nil
}
