package main

import (
	"context"
	log "elicznik-sync/logging"
	"elicznik-sync/tauron"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"time"
)

type TauronDB struct {
	client influxdb2.Client
	org    string
	bucket string
}

func NewTauronDB(url, token, org, bucket string) *TauronDB {
	tdb := TauronDB{}
	tdb.org = org
	tdb.bucket = bucket
	tdb.client = influxdb2.NewClient(url, token)
	return &tdb
}

func (tdb *TauronDB) WriteMeasurements(measurements tauron.DateRangeMeasurements) {
	writeAPI := tdb.client.WriteAPIBlocking(tdb.org, tdb.bucket)

	loc := time.Now().Location()
	for _, date := range measurements.Sort() {
		log.Info("Writing measurements for %s", date)
		ts, _ := time.ParseInLocation(tauron.ELICZNIK_DATE_FORMAT, date, loc)
		for _, h := range measurements[date].Sort() {
			m := (*measurements[date])[h]
			_time := ts.Add(time.Hour * time.Duration(h))
			p := influxdb2.NewPoint("licznik",
				map[string]string{},
				map[string]interface{}{"PowerIn": m.PowerIn, "PowerOut": m.PowerOut},
				_time)
			err := writeAPI.WritePoint(context.Background(), p)
			if err != nil {
				log.Fatal("Error writing to database: %v", err)
			}
		}
	}
}

func (tdb *TauronDB) GetLastMeasurementDate() (time.Time, error) {
	queryAPI := tdb.client.QueryAPI(tdb.org)
	query := `from(bucket: "` + tdb.bucket + `")
		  |> range(start: 0, stop: now())
		  |> filter(fn: (r) => r["_measurement"] == "licznik")
		  |> keep(columns: ["_time"])
		  |> sort(columns: ["_time"], desc: false)
		  |> last(column: "_time")`
	// Get parser flux query result
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return time.Now(), err
	}
	if result.Next() {
		return result.Record().Time().In(time.Now().Location()), nil
	}
	if result.Err() != nil {
		return time.Now(), result.Err()
	}
	return time.Time{}, nil
}
