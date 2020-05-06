package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/naspinall/Hive/pkg/config"

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
		models.WithUsers(cfg.Pepper),
		models.WithMeasurements(),
		models.WithDevices(),
		models.WithAlarms(),
	)

	if err != nil {
		log.Fatal(err)
	}
	defer services.Close()
	services.AutoMigrate()

	usersC := controllers.NewUsers(services.User)
	devicesC := controllers.NewDevices(services.Device)
	measurementsC := controllers.NewMeasurements(services.Measurement)
	alarmsC := controllers.NewAlarms(services.Alarm)
	subscriptionsC := controllers.NewSubscriptions(services.Subscription)
	//auth := middleware.NewJWTAuth(cfg.JWTKey)

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/login", usersC.Login).Methods("POST")

	//Device CRUD, requires JWT Auth(cfg.JWTKey)
	d := api.PathPrefix("/devices").Subrouter()
	//d.Use(auth)
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

	//Measurement CRUD
	m := api.PathPrefix("/measurements").Subrouter()
	//m.Use(auth)
	m.HandleFunc("/{id}/", measurementsC.Delete).Methods("DELETE")
	m.HandleFunc("/{id}/", measurementsC.Get).Methods("GET")

	//Alarm CRUD
	a := api.PathPrefix("/alarms").Subrouter()
	//a.Use(auth)
	a.HandleFunc("/{id}/", alarmsC.Delete).Methods("DELETE")
	a.HandleFunc("/{id}/", alarmsC.Get).Methods("GET")
	a.HandleFunc("/", alarmsC.GetMany).Methods("GET")

	// Subscriptions CRUD
	s := api.PathPrefix("/subscribe").Subrouter()
	//s.Use(auth)
	s.HandleFunc("/", subscriptionsC.Create).Methods("POST")
	s.HandleFunc("/", subscriptionsC.GetMany).Methods("GET")

	//Roles CRUD
	log.Println("Listening on port 3001")
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r)
}
