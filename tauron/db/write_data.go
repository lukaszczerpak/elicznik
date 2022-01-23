package db

import (
	"context"
	log "elicznik/logging"
	"elicznik/tauron/api"
	"elicznik/tauron/balance"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func (db *TauronDB) WriteMeasurements(measurements []api.Measurement) {
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

func (db *TauronDB) WriteBalance(b *balance.TauronMonthlyBalance, complete bool) {
	writeAPI := db.client.WriteAPIBlocking(db.org, db.bucket)
	storageSerialized, err := b.Storage.Marshal()
	if err != nil {
		log.Fatalf("Marshal error: %v", err)
	}
	fields := map[string]interface{}{
		"EnergyPurchased": b.EnergyPurchased,
		"StorageTotal":    b.StorageTotal,
		"Storage":         string(storageSerialized),
	}

	tags := map[string]string{}

	if !complete {
		tags["complete"] = "false"
	}

	measurementName := fmt.Sprintf("%s_balance", db.measurementName)
	p := influxdb2.NewPoint(measurementName, tags, fields, b.Date)
	err = writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		log.Fatalf("Error writing to database: %v", err)
	}
}
