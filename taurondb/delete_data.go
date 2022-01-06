package taurondb

import (
	"context"
	log "elicznik/logging"
	"time"
)

func (tdb *TauronDB) DeleteMeasurements(startTime, stopTime time.Time) {
	deleteAPI := tdb.client.DeleteAPI()
	err := deleteAPI.DeleteWithName(context.Background(), tdb.org, tdb.bucket, startTime, stopTime, "")
	if err != nil {
		log.Errorf("Error when deleting data: %v", err)
	}
}
