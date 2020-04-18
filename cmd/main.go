package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/naspinall/Hive/pkg/controllers"
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
	s := r.PathPrefix("/api").Subrouter()

	//Device CRUD
	s.HandleFunc("/devices", devicesC.GetMany).Methods("GET")
	s.HandleFunc("/devices", devicesC.Create).Methods("POST")
	s.HandleFunc("/devices/{id}/", devicesC.Delete).Methods("DELETE")
	s.HandleFunc("/devices/{id}", devicesC.Get).Methods("GET")

	//User CRUD
	s.HandleFunc("/users", usersC.Create).Methods("POST")
	s.HandleFunc("/users/{id}/", usersC.Delete).Methods("DELETE")
	s.HandleFunc("/users/{id}/", usersC.Get).Methods("GET")

	//Measurement CRUD
	s.HandleFunc("/devices/{id}/measurements", measurementsC.Create).Methods("POST")
	s.HandleFunc("/devices/{id}/measurements", measurementsC.GetByDevice).Methods("GET")
	s.HandleFunc("/measurements/{id}/", measurementsC.Delete).Methods("DELETE")
	s.HandleFunc("/measurements/{id}/", measurementsC.Get).Methods("GET")

	//Alarm CRUD
	s.HandleFunc("/devices/{id}/alarms", alarmsC.Create).Methods("POST")
	s.HandleFunc("/devices/{id}/alarms", alarmsC.GetByDevice).Methods("GET")
	s.HandleFunc("/alarms/{id}/", alarmsC.Delete).Methods("DELETE")
	s.HandleFunc("/alarms/{id}/", alarmsC.Get).Methods("GET")
	s.HandleFunc("/alarms", alarmsC.GetMany).Methods("GET")

	// Subscriptions CRUD
	s.HandleFunc("/devices/{id}/subscribe/", subscriptionsC.Create).Methods("POST")
	s.HandleFunc("/devices/{id}/subscribe/", subscriptionsC.Delete).Methods("DELETE")
	s.HandleFunc("/subscribe/", subscriptionsC.Create).Methods("POST")
	s.HandleFunc("/subscribe/", subscriptionsC.GetMany).Methods("GET")

	//Roles CRUD
	log.Println("Listening on port 3001")
	http.ListenAndServe(":3001", r)
}
