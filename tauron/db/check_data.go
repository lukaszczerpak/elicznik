package db

import (
	"context"
	"fmt"
	"time"
)

func (db *TauronDB) CheckIfDataExists(start, stop time.Time, measurement string) (bool, error) {
	queryAPI := db.client.QueryAPI(db.org)
	query := fmt.Sprintf(`from(bucket: "%v") |> range(start: %v, stop: %v) |> filter(fn: (r) => r["_measurement"] == "%v") |> count(column: "_value")`,
		db.bucket, start.Format(time.RFC3339), stop.Format(time.RFC3339), measurement)
	// Get parser flux query result
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
