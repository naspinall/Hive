package controllers

import (
	"net/http"

	"github.com/naspinall/Hive/pkg/models"
)

func Unauthorized(w http.ResponseWriter, err models.ErrorUnauthorized) {
	http.Error(w, err.Error(), http.StatusUnauthorized)
}
func NotFound(w http.ResponseWriter, err models.ErrorNotFound) {
	http.Error(w, err.Error(), http.StatusNotFound)
}
func InternalServerError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func BadRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func ProcessError(w http.ResponseWriter, err error) {
	if e, ok := err.(models.ErrorNotFound); ok {
		NotFound(w, e)
	} else if e, ok := err.(models.ErrorUnauthorized); ok {
		Unauthorized(w, e)
	} else if e, ok := err.(models.ErrorBadRequest); ok {
		BadRequest(w, e)
	} else {
		InternalServerError(w, err)
	}
	return

}
