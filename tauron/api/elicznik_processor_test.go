package api

import (
	"elicznik/util"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
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

func loadCsvToStruct(filename string, v **ELicznikData) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	eLicznikData, err := ParseTauronCsv(string(data))
	if err != nil {
		log.Fatal(err)
	}

	*v = eLicznikData
}

func loadAndProcessData(filename, startDate, stopDate string) ([]Measurement, error) {
	var data *ELicznikData
	loadCsvToStruct(filename, &data)

	loc, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		log.Fatal(err)
	}

	startTime, _ := time.ParseInLocation(util.DATE_FORMAT, startDate, loc)
	stopTime, _ := time.ParseInLocation(util.DATE_FORMAT, stopDate, loc)

	return processELicznikData(data, startTime, stopTime)
}

func dataComplete(t *testing.T) {
	measurements, err := loadAndProcessData("../../testdata/sample.csv", "2021-10-15", "2021-10-21")
	assert.Len(t, measurements, 7*24, "number of measurements")
	assert.NoError(t, err)
}

func cestToCet(t *testing.T) {
	measurements, err := loadAndProcessData("../../testdata/sample-CEST-to-CET.csv", "2022-10-30", "2022-10-30")
	assert.Len(t, measurements, 25, "number of measurements")
	assert.NoError(t, err)
}

func cetToCest(t *testing.T) {
	measurements, err := loadAndProcessData("../../testdata/sample-CET-to-CEST.csv", "2022-03-27", "2022-03-27")
	assert.Len(t, measurements, 23, "number of measurements")
	assert.NoError(t, err)
}

func arraySizeMismatch(t *testing.T) {
	measurements, err := loadAndProcessData("../../testdata/sample-different-array-sizes.csv", "2021-10-15", "2021-10-21")
	assert.Len(t, measurements, 0)
	assert.Error(t, err)
}

func missingDataOnDay1(t *testing.T) {
	measurements, err := loadAndProcessData("../../testdata/sample-missing-data-on-day1.csv", "2021-10-15", "2021-10-21")
	assert.Len(t, measurements, 0)
	assert.Error(t, err)
}

func missingDataOnDay4(t *testing.T) {
	measurements, err := loadAndProcessData("../../testdata/sample-missing-data-on-day4.csv", "2021-10-15", "2021-10-21")
	assert.Len(t, measurements, 3*24)
	assert.Error(t, err)
}

func missingDataOnDay7(t *testing.T) {
	measurements, err := loadAndProcessData("../../testdata/sample-missing-data-on-day7.csv", "2021-10-15", "2021-10-21")
	assert.Len(t, measurements, 6*24)
	assert.Error(t, err)
}

func missingDay1(t *testing.T) {
	measurements, err := loadAndProcessData("../../testdata/sample-missing-day1.csv", "2021-10-15", "2021-10-21")
	assert.Len(t, measurements, 0)
	assert.Error(t, err)
}

func missingDay4(t *testing.T) {
	measurements, err := loadAndProcessData("../../testdata/sample-missing-day4.csv", "2021-10-15", "2021-10-21")
	assert.Len(t, measurements, 3*24)
	assert.Error(t, err)
}

func missingDay7(t *testing.T) {
	measurements, err := loadAndProcessData("../../testdata/sample-missing-day7.csv", "2021-10-15", "2021-10-21")
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
