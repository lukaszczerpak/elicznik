package api

import (
	log "elicznik/logging"
	"time"
)

type Measurement struct {
	Time     time.Time
	FeedIn   float64
	FromGrid float64
}

func (api *TauronApiClient) GetMeasurements(startTime, stopTime time.Time) ([]Measurement, bool) {
	data, err := api.fetchData(startTime, stopTime)
	if err != nil {
		log.Errorf("Error when fetching data from Tauron API: %v", err)
		return []Measurement{}, false
	}
	measurements, err := processELicznikData(data, startTime, stopTime)
	if err != nil {
		return measurements, false
	}
	return measurements, true
}
