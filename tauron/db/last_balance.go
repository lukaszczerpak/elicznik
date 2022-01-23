package db

import (
	"context"
	"elicznik/tauron/balance"
	"fmt"
)

func (tdb *TauronDB) GetLastCompleteBalance() (*balance.TauronMonthlyBalance, error) {
	queryAPI := tdb.client.QueryAPI(tdb.org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: 0, stop: now())
		|> filter(fn: (r) => 
			r["_measurement"] == "%s_balance" and
			r["complete"] != "false" and
			r["_field"] == "Storage"
		)
		|> keep(columns: ["_time", "_value"])
		|> sort(columns: ["_time"], desc: false)
		|> last(column: "_time")`, tdb.bucket, tdb.measurementName)
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	if result.Next() {
		r := result.Record()
		return balance.NewTauronMonthlyBalance(r.Time().In(tdb.location), r.ValueByKey("_value").(string)), nil
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	return nil, nil
}
