package rest

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/muzz/api/service"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	log       *logrus.Logger
	userConn  service.UserConnector
	validator *validator.Validate
}

func NewHandler(
	log *logrus.Logger,
	userConn service.UserConnector,
) Handler {
	return Handler{
		log:      log,
		userConn: userConn,
		validator: validator.New(
			validator.WithRequiredStructEnabled(),
		),
	}
}

// Healthz godoc
//
//  @Summary      Check service health
//  @Description  Check service health condition
//  @Tags         health
//  @Produce      plain
//  @Success      200  {string}  string  "OK"
//  @Router       /healthz [get]
func (h Handler) Health(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)

	if _, err := writer.Write([]byte(`OK`)); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
