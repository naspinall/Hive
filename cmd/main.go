package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/naspinall/Hive/pkg/config"
	"github.com/naspinall/Hive/pkg/middleware"

	"github.com/gorilla/mux"
	"github.com/naspinall/Hive/pkg/controllers"
	"github.com/naspinall/Hive/pkg/models"
)

func main() {

	cfg := config.LoadConfig()
	dbCfg := cfg.Database

	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(true),
		models.WithSubscriptions(),
		models.WithUsers(cfg.Pepper, cfg.JWTKey),
		models.WithMeasurements(),
		models.WithDevices(),
		models.WithAlarms(),
		models.WithRBAC(),
	)

	if err != nil {
		log.Fatal(err)
	}
	defer services.Close()
	services.DestructiveReset()
	services.AutoMigrate()

	usersC := controllers.NewUsers(services.User, services.RBAC)
	devicesC := controllers.NewDevices(services.Device)
	measurementsC := controllers.NewMeasurements(services.Measurement)
	alarmsC := controllers.NewAlarms(services.Alarm)
	subscriptionsC := controllers.NewSubscriptions(services.Subscription)
	userM := middleware.NewUsersMiddleware(services.User)
	auth := userM.JWTAuth()

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter().StrictSlash(true)

	api.HandleFunc("/login", usersC.Login).Methods("POST")

	//Device CRUD, requires JWT Auth(cfg.JWTKey)
	d := api.PathPrefix("/devices").Subrouter()
	d.Use(auth)
	d.HandleFunc("/", devicesC.GetMany).Methods("GET")
	d.HandleFunc("/", devicesC.Create).Methods("POST")
	d.HandleFunc("/{id}/", devicesC.Delete).Methods("DELETE")
	d.HandleFunc("/{id}", devicesC.Get).Methods("GET")
	d.HandleFunc("/{id}/measurements", measurementsC.Create).Methods("POST")
	d.HandleFunc("/{id}/measurements", measurementsC.GetByDevice).Methods("GET")
	d.HandleFunc("/{id}/alarms", alarmsC.Create).Methods("POST")
	d.HandleFunc("/{id}/alarms", alarmsC.GetByDevice).Methods("GET")
	d.HandleFunc("/{id}/subscribe/", subscriptionsC.Create).Methods("POST")
	d.HandleFunc("/{id}/subscribe/", subscriptionsC.Delete).Methods("DELETE")

	//User CRUD
	u := api.PathPrefix("/users").Subrouter()
	//u.Use(auth)
	u.HandleFunc("/", usersC.Create).Methods("POST")
	u.HandleFunc("/", usersC.GetMany).Methods("GET")
	u.HandleFunc("/{id}/", usersC.Delete).Methods("DELETE")
	u.HandleFunc("/{id}/", usersC.Get).Methods("GET")
	u.HandleFunc("/{id}/roles", usersC.GetRoles).Methods("GET")
	u.HandleFunc("/{id}/roles", usersC.AssignRole).Methods("PUT")

	//Measurement CRUD
	m := api.PathPrefix("/measurements").Subrouter()
	m.Use(auth)
	m.HandleFunc("/{id}/", measurementsC.Delete).Methods("DELETE")
	m.HandleFunc("/{id}/", measurementsC.Get).Methods("GET")

	//Alarm CRUD
	a := api.PathPrefix("/alarms").Subrouter()
	a.Use(auth)
	a.HandleFunc("/{id}/", alarmsC.Delete).Methods("DELETE")
	a.HandleFunc("/{id}/", alarmsC.Create).Methods("POST")
	a.HandleFunc("/{id}/", alarmsC.Get).Methods("GET")
	a.HandleFunc("/", alarmsC.GetMany).Methods("GET")

	// Subscriptions CRUD
	s := api.PathPrefix("/subscribe").Subrouter()
	s.Use(auth)
	s.HandleFunc("/{id}/", subscriptionsC.Create).Methods("POST")
	s.HandleFunc("/{id}/", subscriptionsC.Delete).Methods("DELETE")
	s.HandleFunc("/", subscriptionsC.GetMany).Methods("GET")

	//Roles CRUD
	log.Println(fmt.Sprintf("Listening on port %d", cfg.Port))
	http.ListenAndServe(":3001", r)
}
