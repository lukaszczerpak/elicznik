package db

import (
	"context"
	"fmt"
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

func (tdb *TauronDB) IsMonthComplete(m time.Time) (bool, error) {
	queryAPI := tdb.client.QueryAPI(tdb.org)
	t := time.Date(m.Year(), m.Month(), 1, 23, 0, 0, 0, tdb.location).AddDate(0, 1, -1)
	query := fmt.Sprintf(`from(bucket: "%s")
		  |> range(start: %s, stop: %s)
		  |> filter(fn: (r) => r["_measurement"] == "%s")
		  |> keep(columns: ["_time"])
		  |> sort(columns: ["_time"], desc: false)
		  |> last(column: "_time")`, tdb.bucket, t.Format(time.RFC3339), t.Add(1*time.Second).Format(time.RFC3339), tdb.measurementName)
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return false, err
	}
	if result.Next() {
		return true, nil
	}
	if result.Err() != nil {
		return false, result.Err()
	}
	return false, nil
}
