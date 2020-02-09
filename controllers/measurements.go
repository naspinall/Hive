package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/naspinall/Hive/models"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	measurement.DeviceID = int(id)

	if err := m.ms.Create(&measurement); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (m *Measurements) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := m.ms.Delete(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)

}
func (m *Measurements) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	measurement, err := m.ms.ByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&measurement)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (m *Measurements) GetByDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	measurements, err := m.ms.ByDevice(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&measurements)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
