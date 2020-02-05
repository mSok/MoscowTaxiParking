package main

import (
	"app/controllers"
	"app/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var conf = utils.GetConf()

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/taxiparkings/{id:[0-9]+}", controllers.GetTaxiParkingsByID).Methods("GET")
	router.HandleFunc("/api/load", controllers.LoadTaxiParkings).Methods("GET")
	http.Handle("/", router)

	if conf.Port == "" {
		conf.Port = "8080"
	}
	log.Println("Listen port: ", conf.Port)

	err := http.ListenAndServe(":"+conf.Port, router)
	if err != nil {
		fmt.Print(err)
	}
}
