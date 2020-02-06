package parser

import (
	"app/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func loadFromWeb(url string) ([]byte, error) {
	log.Printf("Load data from url %s", url)
	parkingClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "avtocod")

	res, getErr := parkingClient.Do(req)
	if getErr != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error get data from %s. Status: %s", url, res.Status)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, err
	}
	return body, nil

}

func unpackRawData(data []byte) (*[]models.TaxiParking, error) {
	log.Println("Unpack data from json...")
	parkings := []models.TaxiParking{}
	var raw []json.RawMessage
	jsonErr := json.Unmarshal(data, &raw)
	for _, r := range raw {
		parking := models.TaxiParking{}
		jsonErr = json.Unmarshal(r, &parking)
		if jsonErr != nil {
			return nil, jsonErr
		}
		parking.Raw = r
		parkings = append(parkings, parking)

	}
	if jsonErr != nil {
		return nil, jsonErr
	}
	return &parkings, nil
}

// LoadFromSource load last data from source.
// Source can be a file (file path) or URL.
// Function return number of parsed records and error
func LoadFromSource(source string) (int, error) {
	var data []byte
	var err error
	if strings.HasPrefix(strings.ToLower(source), "http") {
		data, err = loadFromWeb(source)
		if err != nil {
			log.Printf("[Error] %s\n", err)
			return 0, err
		}
	} else {
		log.Printf("Load data from file %s", source)
		data, err = ioutil.ReadFile(source)
		if err != nil {
			log.Printf("[Error] %s\n", err)
			return 0, err
		}
	}
	res, err := unpackRawData(data)
	if err != nil {
		return 0, err
	}
	client := models.GetDB()
	client.BulkInsert(res)
	return len(*res), err
}
