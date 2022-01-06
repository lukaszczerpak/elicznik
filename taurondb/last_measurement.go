package taurondb

import (
	"context"
	"time"
)

func (tdb *TauronDB) GetLastMeasurementDate() (time.Time, error) {
	queryAPI := tdb.client.QueryAPI(tdb.org)
	query := `from(bucket: "` + tdb.bucket + `")
		  |> range(start: 0, stop: now())
		  |> filter(fn: (r) => r["_measurement"] == "` + tdb.measurementName + `")
		  |> keep(columns: ["_time"])
		  |> sort(columns: ["_time"], desc: false)
		  |> last(column: "_time")`
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return time.Now(), err
	}
	if result.Next() {
		yy, mm, dd := result.Record().Time().In(tdb.location).Date()
		return time.Date(yy, mm, dd, 0, 0, 0, 0, tdb.location), nil
	}
	if result.Err() != nil {
		return time.Now(), result.Err()
	}
	return time.Time{}, nil
}
