package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/naspinall/Hive/pkg/controllers"
	"github.com/naspinall/Hive/pkg/middleware"
	"github.com/naspinall/Hive/pkg/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "hive"
	dbname   = "hive"
)

func main() {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	services, err := models.NewServices(connectionString)
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

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/login", usersC.Login).Methods("POST")

	//Device CRUD, requires JWT Auth
	d := api.PathPrefix("/devices").Subrouter()
	d.Use(middleware.Auth)
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
	u.Use(middleware.Auth)
	u.HandleFunc("/", usersC.Create).Methods("POST")
	u.HandleFunc("/", usersC.GetMany).Methods("GET")
	u.HandleFunc("/{id}/", usersC.Delete).Methods("DELETE")
	u.HandleFunc("/{id}/", usersC.Get).Methods("GET")

	//Measurement CRUD
	m := api.PathPrefix("/measurements").Subrouter()
	m.Use(middleware.Auth)
	m.HandleFunc("/{id}/", measurementsC.Delete).Methods("DELETE")
	m.HandleFunc("/{id}/", measurementsC.Get).Methods("GET")

	//Alarm CRUD
	a := api.PathPrefix("/alarms").Subrouter()
	a.Use(middleware.Auth)
	a.HandleFunc("/{id}/", alarmsC.Delete).Methods("DELETE")
	a.HandleFunc("/{id}/", alarmsC.Get).Methods("GET")
	a.HandleFunc("/", alarmsC.GetMany).Methods("GET")

	// Subscriptions CRUD
	s := api.PathPrefix("/subscribe").Subrouter()
	s.Use(middleware.Auth)
	s.HandleFunc("/", subscriptionsC.Create).Methods("POST")
	s.HandleFunc("/", subscriptionsC.GetMany).Methods("GET")

	//Roles CRUD
	log.Println("Listening on port 3001")
	http.ListenAndServe(":3001", r)
}
