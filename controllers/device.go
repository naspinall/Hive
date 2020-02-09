package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/naspinall/Hive/models"
)

type Devices struct {
	ds models.DeviceService
}

func NewDevices(ds models.DeviceService) *Devices {

	return &Devices{
		ds: ds,
	}
}

func (d *Devices) Create(w http.ResponseWriter, r *http.Request) {
	var device models.Device
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := d.ds.Create(&device); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&device)
}

func (d *Devices) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := d.ds.Delete(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)

}
func (d *Devices) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	device, err := d.ds.ByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&device)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
