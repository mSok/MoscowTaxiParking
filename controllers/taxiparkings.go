package controllers

import (
	"app/models"
	"app/parser"
	"app/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetTaxiParkingsByID handler for get ddata by id
var GetTaxiParkingsByID = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	log.Printf("[GET]: %s", params["id"])
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	client := models.GetDB()
	data, err := client.GetTaxiParking(id)
	if data == "" || (err != nil && err.Error() == "redis: nil") {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write([]byte(data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// LoadTaxiParkings handler for load open data
var LoadTaxiParkings = func(w http.ResponseWriter, r *http.Request) {

	var conf = utils.GetConf()
	go func() {
		cnt, err := parser.LoadFromSource(conf.Source)
		if err != nil {
			log.Printf("[Error] parsed data %s\n", err)
		} else {
			log.Printf("[LoadFromSource] loaded %d recs\n", cnt)
		}

	}()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
