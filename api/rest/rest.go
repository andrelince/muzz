package rest

import (
	"net/http"

	_ "github.com/muzz/api/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/muzz/api/rest/middleware"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
)

// @title        Muzz API
// @version      1.0
// @description  This is a API representing a simple dating api system.
// @host         localhost:3000
func NewRest(log *logrus.Logger, router *http.ServeMux, r Handler, auth middleware.AuthMiddleware) error {

	router.HandleFunc("GET /healthz", r.Health)
	router.Handle("GET /swagger/*", httpSwagger.Handler())

	// user
	router.HandleFunc("POST /user/create", r.CreateUser)

	// login
	router.HandleFunc("POST /login", r.Login)

	// swipe
	router.Handle("POST /swipe", auth.Handle(
		http.HandlerFunc(r.Swipe)),
	)

	return nil
}
