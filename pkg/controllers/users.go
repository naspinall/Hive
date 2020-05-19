package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/naspinall/Hive/pkg/models"
)

type LoginResponse struct {
	Token string `json:"token"`
}

type Users struct {
	us    models.UserService
	rbacs models.RBACService
}

func NewUsers(us models.UserService, rbac models.RBACService) *Users {
	return &Users{
		us:    us,
		rbacs: rbac,
	}
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		ProcessError(w, err)
		return
	}

	if err := u.us.Create(&user, r.Context()); err != nil {
		ProcessError(w, err)
		return
	}

	defaultRole := models.Role{
		Alarms:        0,
		Measurements:  0,
		Users:         0,
		Devices:       0,
		Subscriptions: 0,
		UserID:        user.ID,
	}

	if err := u.rbacs.Assign(&defaultRole, r.Context()); err != nil {
		ProcessError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&user)
}

func (u *Users) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ProcessError(w, models.ErrInvalidID)
		return
	}

	if err := u.us.Delete(uint(id), r.Context()); err != nil {
		ProcessError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
func (u *Users) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ProcessError(w, err)
		return
	}

	user, err := u.us.ByID(uint(id), r.Context())
	if err != nil {
		ProcessError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&user)

	if err != nil {
		ProcessError(w, err)
		return
	}
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		ProcessError(w, err)
		return
	}

	fu, err := u.us.Authenticate(user.Email, user.Password, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&fu)
	if err != nil {
		ProcessError(w, err)
		return
	}
}

func (u *Users) GetMany(w http.ResponseWriter, r *http.Request) {
	users, err := u.us.Many(r.Context())
	if err != nil {
		ProcessError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&users)
	if err != nil {
		ProcessError(w, err)
		return
	}

}

func (u *Users) GetRoles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ProcessError(w, models.ErrInvalidID)
		return
	}

	role, err := u.rbacs.ByUserID(uint(id), r.Context())
	if err != nil {
		ProcessError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&role)
	if err != nil {
		ProcessError(w, err)
		return
	}

}

func (u *Users) AssignRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ProcessError(w, models.ErrInvalidID)
		return
	}

	var role models.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		ProcessError(w, err)
		return
	}
	role.UserID = uint(id)

	u.rbacs.Assign(&role, r.Context())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
