package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/muzz/api/rest/definition"
	"github.com/muzz/api/rest/transformer"
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
	v := validator.New(
		validator.WithRequiredStructEnabled(),
	)

	_ = v.RegisterValidation("dob", DOBValidator)

	return Handler{
		log:       log,
		userConn:  userConn,
		validator: v,
	}
}

// Healthz godoc
//
//	@Summary      Check service health
//	@Description  Check service health condition
//	@Tags         health
//	@Produce      plain
//	@Success      200  {string}  string  "OK"
//	@Router       /healthz [get]
func (h Handler) Health(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)

	if _, err := writer.Write([]byte(`OK`)); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateUser godoc
//
// @Summary      Create a user
// @Description  Create a user in the system
// @Tags         user
// @Produce      json
// @Success      200  {object}  definitions.User
// @Router       /user [post]
//
// @Param        user  body  definitions.UserInput  true  "user to create"
func (h Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, errors.New("failed to read request body"))
		return
	}

	var user definition.UserInput
	if err = json.Unmarshal(b, &user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(user); err != nil {
		h.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, err)
		return
	}

	out, err := h.userConn.CreateUser(r.Context(), transformer.FromUserInputDefToEntity(user))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteError(w, err)
		return
	}

	jsonOut, err := json.Marshal(transformer.FromUserEntityToDef(out))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonOut); err != nil {
		WriteError(w, err)
	}
}
