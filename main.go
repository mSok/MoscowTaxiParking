package main

import (
	"app/controllers"
	"app/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.Use(controllers.PrometheusMiddleware)

	router.HandleFunc("/api/taxiparkings/{id:[0-9]+}", controllers.GetTaxiParkingsByID).Methods("GET")

	router.HandleFunc("/api/taxiparkings/localid/{id:[0-9]+}", controllers.GetTaxiParkingsByLocalID).
		Methods("GET").
		Queries("limit", "{limit:[0-9]+}", "offset", "{offset:[0-9]+}")
	// Handle with default limit\offset params
	router.HandleFunc("/api/taxiparkings/localid/{id:[0-9]+}", controllers.GetTaxiParkingsByLocalID).Methods("GET")

	router.HandleFunc("/api/taxiparkings/mode/{mode}", controllers.GetTaxiParkingsByMode).
		Methods("GET").
		Queries("limit", "{limit:[0-9]+}", "offset", "{offset:[0-9]+}")
	// Handle with default limit\offset params
	router.HandleFunc("/api/taxiparkings/mode/{mode}", controllers.GetTaxiParkingsByMode).Methods("GET")

	router.HandleFunc("/api/load", controllers.LoadTaxiParkings).Methods("GET")
	router.HandleFunc("/metrics", controllers.MetricsGetHandler).Methods("GET")
	http.Handle("/", router)
	conf := utils.GetConf()
	log.Println("Listen port: ", conf.Port)

	err := http.ListenAndServe(":"+conf.Port, router)
	if err != nil {
		fmt.Print(err)
	}
}
