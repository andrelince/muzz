package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/muzz/api/pkg/slice"
	"github.com/muzz/api/rest/definition"
	"github.com/muzz/api/rest/middleware"
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

// CreateUser godoc
//
// @Summary      Create a user
// @Description  Create a user in the system
// @Tags         user
// @Produce      json
// @Success      200  {object}  definition.User
// @Router       /user/create [post]
//
// @Param        user  body  definition.UserInput  true  "user to create"
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

// Login godoc
//
// @Summary      Authenticate a user
// @Description  Perform the authentication/login of a user
// @Tags         login
// @Produce      json
// @Success      200  {object}  definition.Token
// @Router       /user [post]
//
// @Param        user  body  definition.LoginInput  true  "credentials to authenticate user"
func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, errors.New("failed to read request body"))
		return
	}

	var login definition.LoginInput
	if err = json.Unmarshal(b, &login); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(login); err != nil {
		h.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, err)
		return
	}

	out, err := h.userConn.Login(r.Context(), login.Email, login.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteError(w, err)
		return
	}

	jsonOut, err := json.Marshal(transformer.FromTokenEntityToDef(out))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonOut); err != nil {
		WriteError(w, err)
	}
}

// Swipe godoc
//
// @Summary      Swipe a user
// @Description  Perform the swipe action on a give user
// @Tags         login
// @Produce      json
// @Success      200  {object}  definition.Match
// @Router       /swipe [post]
//
// @Param        user  body  definition.SwipeInput  true  "swipe data"
func (h Handler) Swipe(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, errors.New("failed to read request body"))
		return
	}

	var swipe definition.SwipeInput
	if err = json.Unmarshal(b, &swipe); err != nil {
		h.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(swipe); err != nil {
		h.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, err)
		return
	}

	var action bool
	if swipe.Preference == "yes" {
		action = true
	}

	out, err := h.userConn.Swipe(r.Context(), userID, swipe.UserID, action)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err)
		WriteError(w, err)
		return
	}

	jsonOut, err := json.Marshal(transformer.FromMatchEntityToDef(out))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonOut); err != nil {
		WriteError(w, err)
	}
}

// Discover godoc
//
// @Summary      Discover relevant profies
// @Description  List profiles of potential match interest
// @Tags         user
// @Produce      json
// @Success      200      {array}  definition.Discovery
// @Param        min_age  query    int     false  "minimum profile age"
// @Param        max_age  query    int     false  "minimum profile age"
// @Param        gender   query    string  false  "M or F"
// @Router       /discover [get]
func (h Handler) Discover(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	params := h.getDiscoverParams(r)

	out, err := h.userConn.Discover(r.Context(), userID, params.Age, params.Gender)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err)
		WriteError(w, err)
		return
	}

	jsonOut, err := json.Marshal(
		slice.Map(out, transformer.FromDiscoveryEntityToDef),
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonOut); err != nil {
		WriteError(w, err)
	}
}

func (h Handler) getDiscoverParams(r *http.Request) discoverParams {
	// Read age parameters
	age := []int{}
	minAgeStr := r.URL.Query().Get("min_age")
	maxAgeStr := r.URL.Query().Get("max_age")
	minAge, err := strconv.Atoi(minAgeStr)
	if err == nil {
		age = append(age, minAge)
	}
	maxAge, err := strconv.Atoi(maxAgeStr)
	if err == nil {
		age = append(age, maxAge)
	}

	// Read gender parameter
	gender := r.URL.Query().Get("gender")

	return discoverParams{
		Gender: gender,
		Age:    age,
	}
}
