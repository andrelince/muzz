package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

func WriteError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
}

func DOBValidator(fl validator.FieldLevel) bool {
	dob := fl.Field().String()
	_, err := time.Parse("2006-01-02", dob)
	return err == nil
}
