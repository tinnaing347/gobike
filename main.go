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
	tag2 := []string{"bike_status"}
	cron1.Every(1).Months(time.Now().Day()).SetTag(tag1).StartImmediately().Do(bike.StationInformationTask)
	cron1.Every(10).Minutes().SetTag(tag2).StartImmediately().Do(bike.StationStatusTask)
	cron1.StartBlocking()
}
