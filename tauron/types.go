package tauron

import (
	"reflect"
	"sort"
)

const USER_AGENT = "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:52.0) Gecko/20100101 Firefox/52.0"
const ELICZNIK_DATE_FORMAT = "2006-01-02"

type ELicznikMeasurement struct {
	Date  string  `json:"Date"`
	Hour  int     `json:"Hour,string"`
	Power float32 `json:"EC,string"`
}

type ELicznikData struct {
	Ok   int `json:"ok"`
	Dane struct {
		PowerOut []ELicznikMeasurement `json:"chart"`
		PowerIn  []ELicznikMeasurement `json:"OZE"`
	} `json:"dane"`
}

type Measurement struct {
	PowerIn  float32 // power sent to operator
	PowerOut float32 // power sent to user
}

type SingleDayMeasurements map[int]*Measurement

type DateRangeMeasurements map[string]*SingleDayMeasurements

func (m *DateRangeMeasurements) Sort() (index []string) {
	for _, k := range reflect.ValueOf(*m).MapKeys() {
		index = append(index, k.String())
	}
	sort.Strings(index)
	return
}

func (m *SingleDayMeasurements) Sort() (index []int) {
	for _, k := range reflect.ValueOf(*m).MapKeys() {
		index = append(index, int(k.Int()))
	}
	sort.Ints(index)
	return
}
