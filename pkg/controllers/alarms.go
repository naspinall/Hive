package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/naspinall/Hive/pkg/models"
)

type Alarms struct {
	as models.AlarmService
}

func NewAlarms(as models.AlarmService) *Alarms {
	return &Alarms{
		as: as,
	}
}

func (a *Alarms) Create(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)

	var alarm models.Alarm
	err = json.NewDecoder(r.Body).Decode(&alarm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	alarm.DeviceID = int(id)

	if err := a.as.Create(&alarm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (a *Alarms) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := a.as.Delete(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)

}
func (a *Alarms) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	alarm, err := a.as.ByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&alarm)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *Alarms) GetByDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	alarms, err := a.as.ByDevice(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&alarms)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
