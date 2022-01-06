package taurondb

import (
	log "elicznik/logging"
	"elicznik/tauronapi"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func (db *TauronDB) WriteMeasurements(measurements []tauronapi.Measurement) {
	writeAPI := db.client.WriteAPI(db.org, db.bucket)
	errorsCh := writeAPI.Errors()
	go func() {
		for err := range errorsCh {
			log.Errorf("Error writing to database: %v", err)
		}
	}()
	for _, m := range measurements {
		fields := map[string]interface{}{
			"FeedIn":   m.FeedIn,
			"FromGrid": m.FromGrid,
		}

		tags := map[string]string{}

		p := influxdb2.NewPoint(db.measurementName, tags, fields, m.Time)
		writeAPI.WritePoint(p)
	}

	writeAPI.Flush()
}
