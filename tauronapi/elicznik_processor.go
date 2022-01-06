package tauronapi

import (
	"elicznik/util"
	"github.com/pkg/errors"
	"time"
)

func processELicznikData(data *ELicznikData, startTime, stopTime time.Time) ([]Measurement, error) {
	stopTime = time.Date(stopTime.Year(), stopTime.Month(), stopTime.Day(), 23, 0, 0, 0, stopTime.Location())

	var measurements = make([]Measurement, 0, len(data.Dane.FromGrid))
	var lastGoodIndex = 0

	if len(data.Dane.FromGrid) != len(data.Dane.FeedIn) {
		return measurements, errors.New("Incomplete data")
	}

	for i, t := 0, startTime; !t.After(stopTime); t, i = t.Add(1*time.Hour), i+1 {
		if t.Hour() == 0 {
			lastGoodIndex = i
		}

		if i >= len(data.Dane.FeedIn) {
			return measurements[:lastGoodIndex], errors.New("Incomplete data")
		}

		feedIn := data.Dane.FeedIn[i]
		fromGrid := data.Dane.FromGrid[i]

		if feedIn.Date != fromGrid.Date || feedIn.Hour != fromGrid.Hour ||
			feedIn.Date != t.Format(util.DATE_FORMAT) || feedIn.Hour-1 != t.Hour() {
			return measurements[:lastGoodIndex], errors.New("Incomplete data")
		}

		measurements = append(measurements, Measurement{t, feedIn.Power * 1000, fromGrid.Power * 1000})
	}

	return measurements, nil
}
