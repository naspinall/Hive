package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/naspinall/Hive/pkg/models"
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
		ProcessError(w, err)
		return
	}

	if err := d.ds.Create(&device, r.Context()); err != nil {
		ProcessError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&device)
}

func (d *Devices) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ProcessError(w, models.ErrInvalidID)
		return
	}

	if err := d.ds.Delete(uint(id), r.Context()); err != nil {
		ProcessError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
func (d *Devices) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ProcessError(w, err)
		return
	}

	device, err := d.ds.ByID(uint(id), r.Context())
	if err != nil {
		ProcessError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&device)

	if err != nil {
		ProcessError(w, err)
		return
	}
}

func (d *Devices) GetMany(w http.ResponseWriter, r *http.Request) {
	var err error
	var count int64 = 100
	q := r.URL.Query()
	cq, ok := q["count"]
	if ok {
		count, err = strconv.ParseInt(cq[0], 10, 64)
		if err != nil {
			ProcessError(w, models.ErrInvalidID)
			return
		}
	}
	devices, err := d.ds.Many(int(count), r.Context())
	if err != nil {
		ProcessError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&devices)

	if err != nil {
		ProcessError(w, err)
		return
	}
}

func (d *Devices) GetByName(w http.ResponseWriter, r *http.Request) {
	var name string
	q := r.URL.Query()
	cq, ok := q["name"]
	if ok {
		name = cq[0]
	}
	if ok != true {
		http.Error(w, "Name required for search", http.StatusBadRequest)
		return
	}

	devices, err := d.ds.SearchByName(name, r.Context())
	if err != nil {
		ProcessError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&devices)

	if err != nil {
		ProcessError(w, err)
		return
	}
}
