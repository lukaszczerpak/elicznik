package db

import (
	"context"
	log "elicznik/logging"
	"fmt"
	"time"
)

func (tdb *TauronDB) delete(startTime, stopTime time.Time, measurement string) {
	deleteAPI := tdb.client.DeleteAPI()
	err := deleteAPI.DeleteWithName(context.Background(), tdb.org, tdb.bucket, startTime, stopTime,
		fmt.Sprintf(`_measurement="%s"`, measurement))
	if err != nil {
		log.Errorf("Error when deleting data: %v", err)
	}
}

func (tdb *TauronDB) DeleteMeasurements(startTime, stopTime time.Time) {
	tdb.delete(startTime, stopTime, tdb.measurementName)
}

func (tdb *TauronDB) DeleteBalance(startTime, stopTime time.Time) {
	tdb.delete(startTime, stopTime, fmt.Sprintf("%s_balance", tdb.measurementName))
}
