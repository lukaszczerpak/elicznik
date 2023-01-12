package db

import (
	"context"
	"fmt"
	"time"
)

func (tdb *TauronDB) GetMonthStats(m time.Time) (float64, float64, error) {
	queryAPI := tdb.client.QueryAPI(tdb.org)
	start := time.Date(m.Year(), m.Month(), 1, 0, 0, 0, 0, tdb.location)
	stop := start.AddDate(0, 1, 0).Add(-1 * time.Second)
	query := fmt.Sprintf(`
		import "influxdata/influxdb/schema"
		import "date"

		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "%s" and (r["_field"] == "FeedIn" or r["_field"] == "FromGrid"))
		|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
		|> reduce(
			  identity: {FeedIn_sum: 0.0, FromGrid_sum: 0.0},
			  fn: (r, accumulator) => ({
				FeedIn_sum: r.FeedIn + accumulator.FeedIn_sum,
				FromGrid_sum: r.FromGrid + accumulator.FromGrid_sum
			  })
			)
		|> rename(columns: {FeedIn_sum: "FeedIn", FromGrid_sum: "FromGrid"})
		|> keep(columns: ["FeedIn", "FromGrid"])`,
		tdb.bucket, start.Format(time.RFC3339), stop.Format(time.RFC3339), tdb.measurementName)
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return -1, -1, err
	}
	if result.Next() {
		r := result.Record()
		return r.ValueByKey("FromGrid").(float64), r.ValueByKey("FeedIn").(float64), nil
	}
	if result.Err() != nil {
		return -1, -1, result.Err()
	}
	return 0, 0, nil
}
