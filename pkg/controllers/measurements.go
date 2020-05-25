package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/naspinall/Hive/pkg/models"
)

type Measurements struct {
	ms models.MeasurementService
}

func NewMeasurements(ms models.MeasurementService) *Measurements {
	return &Measurements{
		ms: ms,
	}
}

func (m *Measurements) Create(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)

	var measurement models.Measurement
	err = json.NewDecoder(r.Body).Decode(&measurement)
	if err != nil {
		ProcessError(w, err)
		return
	}
	measurement.DeviceID = uint(id)

	if err := m.ms.Create(&measurement, r.Context()); err != nil {
		ProcessError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (m *Measurements) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ProcessError(w, models.ErrInvalidID)
		return
	}

	if err := m.ms.Delete(uint(id), r.Context()); err != nil {
		ProcessError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
func (m *Measurements) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ProcessError(w, models.ErrInvalidID)
		return
	}

	measurement, err := m.ms.ByID(uint(id), r.Context())
	if err != nil {
		ProcessError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&measurement)

	if err != nil {
		ProcessError(w, err)
		return
	}
}

func (m *Measurements) GetByDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ProcessError(w, models.ErrInvalidID)
		return
	}

	measurements, err := m.ms.ByDevice(uint(id), r.Context())
	if err != nil {
		ProcessError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&measurements)

	if err != nil {
		ProcessError(w, err)
		return
	}
}
