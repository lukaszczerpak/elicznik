package tauronapi

import (
	"elicznik/util"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

func getFunctionName(temp interface{}) string {
	strs := strings.Split(runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name(), ".")
	return strs[len(strs)-1]
}

func loadAndProcessData(filename, startDate, stopDate string) ([]Measurement, error) {
	var data *ELicznikData
	loadJsonToStruct(filename, &data)
	loc, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		log.Fatal(err)
	}
	startTime, _ := time.ParseInLocation(util.DATE_FORMAT, startDate, loc)
	stopTime, _ := time.ParseInLocation(util.DATE_FORMAT, stopDate, loc)
	return processELicznikData(data, startTime, stopTime)
}

func dataComplete(t *testing.T) {
	measurements, err := loadAndProcessData("../testdata/sample.json", "2021-10-15", "2021-10-21")
	assert.Len(t, measurements, 7*24, "number of measurements")
	assert.NoError(t, err)
}

func cestToCet(t *testing.T) {
	measurements, err := loadAndProcessData("../testdata/sample-CEST-to-CET.json", "2021-10-31", "2021-10-31")
	assert.Len(t, measurements, 25, "number of measurements")
	assert.NoError(t, err)
}

func cetToCest(t *testing.T) {
	measurements, err := loadAndProcessData("../testdata/sample-CET-to-CEST.json", "2021-03-28", "2021-03-28")
	assert.Len(t, measurements, 23, "number of measurements")
	assert.NoError(t, err)
}

func arraySizeMismatch(t *testing.T) {
	measurements, err := loadAndProcessData("../testdata/sample-different-array-sizes.json", "2021-10-15", "2021-10-21")
	assert.Len(t, measurements, 0)
	assert.Error(t, err)
}

func missingDataOnDay1(t *testing.T) {
	measurements, err := loadAndProcessData("../testdata/sample-missing-data-on-day1.json", "2021-10-15", "2021-10-21")
	assert.Len(t, measurements, 0)
	assert.Error(t, err)
}

func missingDataOnDay4(t *testing.T) {
	measurements, err := loadAndProcessData("../testdata/sample-missing-data-on-day4.json", "2021-10-15", "2021-10-21")
	assert.Len(t, measurements, 3*24)
	assert.Error(t, err)
}

func missingDataOnDay7(t *testing.T) {
	measurements, err := loadAndProcessData("../testdata/sample-missing-data-on-day7.json", "2021-10-15", "2021-10-21")
	assert.Len(t, measurements, 6*24)
	assert.Error(t, err)
}

func missingDay1(t *testing.T) {
	measurements, err := loadAndProcessData("../testdata/sample-missing-day1.json", "2021-10-15", "2021-10-21")
	assert.Len(t, measurements, 0)
	assert.Error(t, err)
}

func missingDay4(t *testing.T) {
	measurements, err := loadAndProcessData("../testdata/sample-missing-day4.json", "2021-10-15", "2021-10-21")
	assert.Len(t, measurements, 3*24)
	assert.Error(t, err)
}

func missingDay7(t *testing.T) {
	measurements, err := loadAndProcessData("../testdata/sample-missing-day7.json", "2021-10-15", "2021-10-21")
	assert.Len(t, measurements, 6*24)
	assert.Error(t, err)
}

func TestProcessELicznikData(t *testing.T) {
	var tests = []func(t *testing.T){
		dataComplete,
		cestToCet,
		cetToCest,
		arraySizeMismatch,
		missingDataOnDay1,
		missingDataOnDay4,
		missingDataOnDay7,
		missingDay1,
		missingDay4,
		missingDay7,
	}

	for _, fn := range tests {
		t.Run(getFunctionName(fn), fn)
	}
}
