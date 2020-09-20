package bike

import (
	"strconv"
	"time"
)

type StationInformation struct {
	Data        StationInformationDict `json:"data"`
	LastUpdated Time                   `json:"last_updated"`
	TTL         int                    `json:"ttl"`
}

//unwrap one more station layer
type StationInformationDict struct {
	Stations []Station `json:"stations"`
}

type Station struct {
	LegacyID                    string   `json:"legacy_id"`
	StationID                   string   `json:"station_id"`
	Name                        string   `json:"name"`
	Capacity                    int      `json:"capacity"`
	ElectricBikeSurchargeWaiver bool     `json:"electric_bike_surcharge_waiver"`
	Lat                         float32  `json:"lat"`
	Lon                         float32  `json:"lon"`
	RegionID                    string   `json:"region_id"`
	StationType                 string   `json:"station_type"`
	RentalMethods               []string `json:"rental_methods"`
	HasKiosk                    bool     `json:"has_kiosk"`
}

type StationStatus struct {
	Data        StationStatusDict `json:"data"`
	LastUpdated Time              `json:"last_updated"`
	TTL         int               `json:"ttl"`
}

//unwrap one more station layer
type StationStatusDict struct {
	Stations []BikeStatusAtStation `json:"stations"`
}

type BikeStatusAtStation struct {
	IsReturning            int    `json:"is_returning"`
	NumBikesDisabled       int    `json:"num_bikes_disabled"`
	EightdHasAvailableKeys bool   `json:"eightd_has_available_keys"`
	NumEbikesAvailable     int    `json:"num_ebikes_available"`
	StationStat            string `json:"station_status"`
	IsRenting              int    `json:"is_renting"`
	LastReported           Time   `json:"last_reported"`
	NumDocksAvaliable      int    `json:"num_docks_available"`
	IsInstalled            int    `json:"is_installed"`
	StationID              string `json:"station_id"`
	LegacyID               string `json:"legacy_id"`
	NumDocksDisabled       int    `json:"num_docks_disabled"`
}

//discoverd the following from https://www.yellowduck.be/posts/handling-unix-timestamps-in-json/ just like colombus

// Time defines a timestamp encoded as epoch seconds in JSON
type Time time.Time

// UnmarshalJSON is used to convert the timestamp from JSON
func (t *Time) UnmarshalJSON(s []byte) (err error) {
	r := string(s)
	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(q, 0)
	return nil
}

// Unix returns t as a Unix time, the number of seconds elapsed
// since January 1, 1970 UTC. The result does not depend on the
// location associated with t.
func (t Time) Unix() int64 {
	return time.Time(t).Unix()
}

// Time returns the JSON time as a time.Time instance in UTC
func (t Time) Time() time.Time {
	return time.Time(t).UTC()
}

// String returns t as a formatted string
func (t Time) String() string {
	return t.Time().String()
}

func (b *BikeStatusAtStation) TagField() (map[string]string, map[string]interface{}, time.Time) {

	tags := map[string]string{"station_id": b.StationID, "legacy_id": b.LegacyID}
	fields := map[string]interface{}{
		"is_returning":              b.IsReturning,
		"num_bikes_disabled":        b.NumBikesDisabled,
		"eightd_has_available_keys": b.EightdHasAvailableKeys,
		"num_ebikes_available":      b.NumEbikesAvailable,
		"station_status":            b.StationStat,
		"is_renting":                b.IsRenting,
		"num_docks_available":       b.NumDocksAvaliable,
		"is_installed":              b.IsInstalled,
		"num_docks_disabled":        b.NumDocksDisabled,
	}

	return tags, fields, b.LastReported.Time()
}

func (s *Station) TagField(lastUpdated Time) (map[string]string, map[string]interface{}, time.Time) {

	tags := map[string]string{"station_id": s.StationID, "legacy_id": s.LegacyID, "region_id": s.RegionID}
	fields := map[string]interface{}{
		"name":                           s.Name,
		"capacity":                       s.Capacity,
		"electric_bike_surcharge_waiver": s.ElectricBikeSurchargeWaiver,
		"lat":                            s.Lat,
		"lon":                            s.Lon,
		"station_type":                   s.StationType,
		"has_kiosk":                      s.HasKiosk,
	}

	return tags, fields, lastUpdated.Time()
}
