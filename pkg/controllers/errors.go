package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/naspinall/Hive/pkg/models"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func Unauthorized(w http.ResponseWriter, err models.ErrorUnauthorized) {
	er := ErrorResponse{Message: err.Error(), Status: http.StatusUnauthorized}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(er)
}
func NotFound(w http.ResponseWriter, err models.ErrorNotFound) {
	er := ErrorResponse{Message: err.Error(), Status: http.StatusNotFound}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(er)
}
func InternalServerError(w http.ResponseWriter, err error) {
	er := ErrorResponse{Message: err.Error(), Status: http.StatusInternalServerError}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(er)
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
