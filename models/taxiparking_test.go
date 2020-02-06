package models

import (
	"encoding/json"
	"log"
	"strconv"

	"os"
	"testing"

	"github.com/go-redis/redis"
)

var raw = json.RawMessage(`{
"global_id": 1704691,
"system_object_id": "161",
"ID": 161,
"Name": "Парковка такси по адресу Карачаровское шоссе, дом 15",
"AdmArea": "Юго-Восточный административный округ",
"District": "Нижегородский район",
"Address": "Карачаровское шоссе, дом 15",
"Longitude_WGS84": "37.7630192041397",
"Latitude_WGS84": "55.7356914963956",
"CarCapacity": 4,
"Mode": "круглосуточно",
"ID_en": 161,
"Name_en": "Taxi parking at Karacharovskoe shosse, house 15",
"AdmArea_en": "Yugo-Vostochny'j administrativny'j okrug",
"District_en": "Nizhegorodskij rajon",
"Address_en": "Karacharovskoe shosse, dom 15",
"Longitude_WGS84_en": "37.7630192041397",
"Latitude_WGS84_en": "55.7356914963956",
"CarCapacity_en": 4,
"Mode_en": "24-hours"
}`)

var mockdata = []TaxiParking{
	TaxiParking{GlobalID: 1704691,
		ID:     161,
		ModeEN: "24-hours",
		Raw:    raw,
	}}

// Create a test connection with Redis. Be careful, all data in the database will be reset.
// Number of db for test, pass in env testdb=...
func connectToTest() *DBClient {
	addr := os.Getenv("redis")
	password := os.Getenv("password")
	var (
		dbNum int
		err   error
	)
	if os.Getenv("testdb") == "" {
		dbNum = 0
	} else {
		dbNum, err = strconv.Atoi(os.Getenv("testdb"))
		if err != nil {
			log.Fatal(err)
		}

	}

	dbredis := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbNum,
	})
	dbredis.FlushDB()
	return &DBClient{db: dbredis}
}

func TestBulkInsert(t *testing.T) {
	client := connectToTest()
	client.BulkInsert(&mockdata)

	prefixKey, err := client.getActualPrefix()
	if err != nil {
		t.Errorf("BulkInser. getActualPrefix error %s", err)
	}
	if prefixKey == "" {
		t.Errorf("BulkInsert not working. There is nothing in \"parking:keys\"")
	}

	res, err := client.GetTaxiParking(1704691)
	if err != nil {
		t.Errorf("GetTaxiParking not working. Error %s", err)
	}
	if res == "" {
		t.Errorf("BulkInsert not working. GetTaxiParking return nothing")
	}
	client.db.FlushDB()
}

func TestGetParkingTaxiById(t *testing.T) {
	client := connectToTest()
	client.BulkInsert(&mockdata)
	res, _ := client.GetTaxiParking(1704691)
	if res != string(raw) {
		t.Errorf("Error compare data GetTaxiParking with the source")
	}
	// first 1, in range with leght 1
	resArr, _ := client.GetTaxiParkingByID(161, 1, 0)
	if len(resArr) != 1 {
		t.Errorf("Len of result GetTaxiParkingByID is not equal 1")
	}
	// first 10, in range with leght 1
	resArr, _ = client.GetTaxiParkingByID(161, 10, 0)
	if len(resArr) != 1 {
		t.Errorf("Len of result GetTaxiParkingByID is not equal 1")
	}
	// first 1 with offset 10, in range with leght 1
	resArr, _ = client.GetTaxiParkingByID(161, 1, 10)
	if len(resArr) != 0 {
		t.Errorf("Offset wrong working in GetTaxiParkingByID")
	}

	client.db.FlushDB()
	res, _ = client.GetTaxiParking(1704691)
	if res != "" {
		t.Errorf("Incorrect data in emty DB")
	}
}

func TestGetParkingTaxiByMode(t *testing.T) {
	client := connectToTest()
	client.BulkInsert(&mockdata)
	// first 1, in range with leght 1
	resArr, _ := client.GetTaxiParkingByMode("24-hours", 1, 0)
	if len(resArr) != 1 {
		t.Errorf("Len of result GetTaxiParkingByID is wrong")
	}
	// first 1 with offset 10, in range with leght 1
	resArr, _ = client.GetTaxiParkingByMode("24-hours", 1, 10)
	if len(resArr) != 0 {
		t.Errorf("Len of result GetTaxiParkingByID is wrong")
	}

}
