package main

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/tinnaing347/gobike/bike"
	"github.com/tinnaing347/gobike/models"
)

func main() {

	models.CreateClient()
	defer models.DB.Close()

	cron1 := gocron.NewScheduler(time.UTC)
	tag1 := []string{"station_information"}
	cron1.Every(1).Minutes().SetTag(tag1).Do(bike.StationInformationTask)
	// cron1.Every(2).Minutes([]string{"bike_status"}).Lock().Do(bike.StationStatusTask)
	cron1.StartBlocking()
}
