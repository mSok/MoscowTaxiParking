package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

// TaxiParking  - description infomation about taxi parking in Moscow.
type TaxiParking struct {
	GlobalID int    `json:"global_id"`
	ID       int    `json:"ID"`
	ModeEN   string `json:"Mode_en"`
	Raw      json.RawMessage
}

// Get actual prefix keys in database
func (d *DBClient) getActualPrefix() (string, error) {
	client := d.db
	str, err := client.LRange("parking:keys", -1, -1).Result()
	if err != nil {
		return "", err
	}
	if len(str) == 0 {
		return "", nil
	}
	return str[0], nil
}

// Get json data from a database by list of global ids
func (d *DBClient) getRangeByIds(ids []string) ([]string, error) {
	parkings := []string{}
	for _, gid := range ids {
		numGid, err := strconv.Atoi(gid)
		if err != nil {
			log.Printf("Error convert global_id %s", gid)
			continue
		}

		p, err := d.GetTaxiParking(numGid)
		if err != nil {
			log.Printf("Error get data by id %s", err)
			continue
		}
		if p != "" {
			parkings = append(parkings, p)
		}
	}
	return parkings, nil
}

// BulkInsert - bulk insert json data into the db and create indexes to search
func (d *DBClient) BulkInsert(data *[]TaxiParking) {
	// create prefix for actual data
	prefixKey := strconv.FormatInt(time.Now().Unix(), 10)
	client := d.db
	client.RPush("parking:keys", prefixKey)

	for _, d := range *data {
		client.Set(fmt.Sprintf("parking:globalid:%s:%d", prefixKey, d.GlobalID), string(d.Raw), 0)
		// Save index for serach by Id
		if d.ID != 0 {
			client.RPush(fmt.Sprintf("parking:id:%s:%d", prefixKey, d.ID), d.GlobalID)
		}
		// Save index for serach by Mode
		if d.ModeEN != "" {
			client.RPush(fmt.Sprintf("parking:mode:%s:%s", prefixKey, d.ModeEN), d.GlobalID)
		}
	}
	log.Printf("Actula prefix %s", prefixKey)
	// TODO Need to remove outdated data
}

// GetTaxiParking return json string by global_id
func (d *DBClient) GetTaxiParking(GlobalID int) (string, error) {
	prefixKey, err := d.getActualPrefix()
	if err != nil {
		return "", err
	}
	if prefixKey == "" {
		return "", nil
	}
	client := d.db
	return client.Get(fmt.Sprintf("parking:globalid:%s:%d", prefixKey, GlobalID)).Result()
}

// GetTaxiParking return json string by ID
// Since it is not clear whether the ID is unique, return an array of json data by ID
func (d *DBClient) GetTaxiParkingByID(ID int, limit int, offset int) ([]string, error) {
	prefixKey, err := d.getActualPrefix()
	if err != nil {
		return nil, err
	}
	client := d.db
	res, err := client.LRange(fmt.Sprintf("parking:id:%s:%d", prefixKey, ID), int64(offset), int64(limit)).Result()
	if err != nil {
		return nil, err
	}
	return d.getRangeByIds(res)
}

// GetTaxiParkingByMode return an array of json data by Mode
func (d *DBClient) GetTaxiParkingByMode(Mode string, limit int, offset int) ([]string, error) {
	prefixKey, err := d.getActualPrefix()
	if err != nil {
		return nil, err
	}
	client := d.db
	res, err := client.LRange(fmt.Sprintf("parking:mode:%s:%s", prefixKey, Mode), int64(offset), int64(limit)).Result()
	if err != nil {
		return nil, err
	}
	return d.getRangeByIds(res)
}
