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
	router.HandleFunc("/api/taxiparkings/{id:[0-9]+}", controllers.GetTaxiParkingsByID).Methods("GET")
	router.HandleFunc("/api/load", controllers.LoadTaxiParkings).Methods("GET")
	http.Handle("/", router)
	conf := utils.GetConf()
	log.Println("Listen port: ", conf.Port)

	err := http.ListenAndServe(":"+conf.Port, router)
	if err != nil {
		fmt.Print(err)
	}
}
