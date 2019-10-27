package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/naspinall/Hive/controllers"
	"github.com/naspinall/Hive/models"
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

	r := mux.NewRouter()

	//Device CRUD
	r.HandleFunc("/device", devicesC.Create).Methods("POST")
	r.HandleFunc("/device/{id}/", devicesC.Delete).Methods("DELETE")
	r.HandleFunc("/device/{id}/", devicesC.Get).Methods("GET")

	//User CRUD
	r.HandleFunc("/user", usersC.Create).Methods("POST")
	r.HandleFunc("/user/{id}/", usersC.Delete).Methods("DELETE")
	r.HandleFunc("/user/{id}/", usersC.Get).Methods("GET")

	//Measurement CRUD
	r.HandleFunc("/device/{id}/measurement", measurementsC.Create).Methods("POST")
	r.HandleFunc("/device/{id}/measurement", measurementsC.GetByDevice).Methods("GET")
	r.HandleFunc("/measurement/{id}/", measurementsC.Delete).Methods("DELETE")
	r.HandleFunc("/measurement/{id}/", measurementsC.Get).Methods("GET")

	//Alarm CRUD
	r.HandleFunc("/device/{id}/alarm", alarmsC.Create).Methods("POST")
	r.HandleFunc("/device/{id}/alarm", alarmsC.GetByDevice).Methods("GET")
	r.HandleFunc("/alarm/{id}/", alarmsC.Delete).Methods("DELETE")
	r.HandleFunc("/alarm/{id}/", alarmsC.Get).Methods("GET")

	//Roles CRUD

	http.ListenAndServe(":3000", r)
}
