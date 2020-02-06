package parser

import (
	"testing"
)

const mockdata = `[{
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
}]`

func TestLoadFromSource(t *testing.T) {
	_, err := LoadFromSource("")
	if err == nil {
		t.Error("No Error with empty source")
	}
}

func TestUnpack(t *testing.T) {
	recs, err := unpackRawData([]byte(mockdata))
	if err != nil {
		t.Errorf("Error unmarshal JSON %s", err)
	}
	if len((*recs)) != 1 {
		t.Errorf("Error unmarshal JSON, lenght %d not equal 1", len((*recs)))
	}
	if (*recs)[0].GlobalID != 1704691 {
		t.Errorf("Error unmarshal JSON, GlobalID %d not equal source (1704691)", (*recs)[0].GlobalID)
	}

}
