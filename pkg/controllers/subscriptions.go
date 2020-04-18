package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/naspinall/Hive/pkg/models"
)

type Subscriptions struct {
	ss models.SubscriptionService
}

func NewSubscriptions(ss models.SubscriptionService) *Subscriptions {
	return &Subscriptions{
		ss: ss,
	}
}

func (s *Subscriptions) Create(w http.ResponseWriter, r *http.Request) {

	var id int64
	var err error

	vars := mux.Vars(r)
	if vars["id"] == "*" {
		id = -1
	} else {
		id, err = strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, "Bad ID", http.StatusBadRequest)
		}
	}

	var subscription models.Subscription
	err = json.NewDecoder(r.Body).Decode(&subscription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	subscription.DeviceID = int(id)

	if err = s.ss.Create(&subscription); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&subscription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Subscriptions) GetMany(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := s.ss.Many()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(subscriptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Subscriptions) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := s.ss.Delete(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)

}
