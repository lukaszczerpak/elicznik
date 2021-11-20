package tauron

import (
	log "elicznik-sync/logging"
	"time"
)

type DataFetcher func(dateFrom time.Time, dateTo time.Time) (*ELicznikData, error)

type ELicznik struct {
	fetcher DataFetcher
}

func initArray() *SingleDayMeasurements {
	arr := make(SingleDayMeasurements)
	for i := 1; i <= 24; i++ {
		arr[i] = &Measurement{-1, -1}
	}
	return &arr
}

func isMeasurementDayComplete(measurements *SingleDayMeasurements) bool {
	for _, m := range *measurements {
		if m.PowerOut == -1 || m.PowerIn == -1 {
			return false
		}
	}
	return true
}

func processData(data *ELicznikData, dateFrom time.Time, dateTo time.Time) (DateRangeMeasurements, bool) {
	// prepare structures
	measurementsByDay := make(DateRangeMeasurements)
	dateStart := dateFrom.Truncate(24)
	dateEnd := dateTo.Truncate(24)
	for day := 0; dateStart.AddDate(0, 0, day).Before(dateEnd.AddDate(0, 0, 1)); day++ {
		date := dateStart.AddDate(0, 0, day).Format(ELICZNIK_DATE_FORMAT)
		measurementsByDay[date] = initArray()
	}

	// iterate over measurements
	for i := range data.Dane.PowerIn {
		e := data.Dane.PowerIn[i]
		(*measurementsByDay[e.Date])[e.Hour].PowerIn = e.Power

		e = data.Dane.PowerOut[i]
		(*measurementsByDay[e.Date])[e.Hour].PowerOut = e.Power
	}

	deleteRemaining := false
	for day := 0; dateStart.AddDate(0, 0, day).Before(dateEnd.AddDate(0, 0, 1)); day++ {
		date := dateStart.AddDate(0, 0, day).Format(ELICZNIK_DATE_FORMAT)
		if deleteRemaining || !isMeasurementDayComplete(measurementsByDay[date]) {
			log.Warn("Incomplete day: %s, deleting from the dataset", date)
			delete(measurementsByDay, date)
			deleteRemaining = true
			continue
		}
	}

	return measurementsByDay, deleteRemaining
}

func NewELicznik(username, password, smartNr string) *ELicznik {
	fetcher := NewELicznikFetcher(username, password, smartNr)
	elicznik := ELicznik{fetcher: fetcher.fetchData}
	return &elicznik
}

func (elicznik *ELicznik) GetData(dateFrom time.Time, dateTo time.Time) (DateRangeMeasurements, bool, error) {
	data, err := elicznik.fetcher(dateFrom, dateTo)
	if err != nil {
		return nil, false, err
	}
	processedData, incomplete := processData(data, dateFrom, dateTo)
	return processedData, incomplete, nil
}
