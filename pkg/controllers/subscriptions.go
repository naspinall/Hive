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
			ProcessError(w, models.ErrInvalidID)
			return
		}
	}

	var subscription models.Subscription
	err = json.NewDecoder(r.Body).Decode(&subscription)
	if err != nil {
		ProcessError(w, err)
		return
	}

	subscription.DeviceID = uint(id)

	if err = s.ss.Create(&subscription, r.Context()); err != nil {
		ProcessError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&subscription)
	if err != nil {
		ProcessError(w, err)
		return
	}
}

func (s *Subscriptions) GetMany(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := s.ss.Many(r.Context())
	if err != nil {
		ProcessError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(subscriptions)
	if err != nil {
		ProcessError(w, err)
		return
	}
}

func (s *Subscriptions) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ProcessError(w, models.ErrInvalidID)
		return
	}

	if err := s.ss.Delete(uint(id), r.Context()); err != nil {
		ProcessError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
