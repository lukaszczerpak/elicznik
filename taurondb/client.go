package taurondb

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"time"
)

type TauronDB struct {
	client          influxdb2.Client
	org             string
	bucket          string
	location        *time.Location
	measurementName string
}

func New(loc *time.Location, url, token, org, bucket, measurementName string) *TauronDB {
	db := TauronDB{}
	db.location = loc
	db.org = org
	db.bucket = bucket
	db.client = influxdb2.NewClient(url, token)
	db.measurementName = measurementName
	return &db
}

func (db *TauronDB) Close() {
	db.client.Close()
}
