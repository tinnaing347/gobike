package bike

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/tinnaing347/gobike/models"
)

// get station_information response from citibike
func StationInformationResponse() *StationInformation {
	resp, err := http.Get("https://gbfs.citibikenyc.com/gbfs/fr/station_information.json")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var sInfo StationInformation
	if err = json.NewDecoder(resp.Body).Decode(&sInfo); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	return &sInfo
}

// get station_status response from citibike
func StationStatusResponse() *StationStatus {
	resp, err := http.Get("https://gbfs.citibikenyc.com/gbfs/en/station_status.json")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var sStat StationStatus
	if err = json.NewDecoder(resp.Body).Decode(&sStat); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	return &sStat
}

func PopulateStations(s StationInformation) {

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "bikedb",
		Precision: "s",
	})

	if err != nil {
		log.Fatal(err)
	}

	station_information := s.Data.Stations
	lastUpdated := s.LastUpdated

	for i := 0; i < len(station_information); i++ {
		tags, fields, time_ := station_information[i].TagField(lastUpdated)
		pt, err := client.NewPoint("stations", tags, fields, time_)
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}
	if err := models.DB.Write(bp); err != nil {
		log.Fatal(err)
	}
	log.Println("sucessuflly written data for stations")
}

func PopulateBikeStatusAtStation(s StationStatus) {

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "bikedb",
		Precision: "s",
	})

	if err != nil {
		log.Fatal(err)
	}

	bikeAtStation := s.Data.Stations

	for i := 0; i < len(bikeAtStation); i++ {
		tags, fields, time_ := bikeAtStation[i].TagField()
		pt, err := client.NewPoint("bike_status_at_station", tags, fields, time_)
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}
	if err := models.DB.Write(bp); err != nil {
		log.Fatal(err)
	}
	log.Println("sucessuflly written data for bike_status_at_station")
}

func StationStatusTask() {
	s := StationStatusResponse()
	PopulateBikeStatusAtStation(*s)
}

func StationInformationTask() {
	s := StationInformationResponse()
	PopulateStations(*s)
}
