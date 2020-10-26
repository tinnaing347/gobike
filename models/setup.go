package models

import (
	"log"

	client "github.com/influxdata/influxdb1-client/v2"
)

var DB client.Client

func CreateClient(db_address string) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: db_address,
	})

	if err != nil {
		log.Fatal(err)
	}

	DB = c
}
