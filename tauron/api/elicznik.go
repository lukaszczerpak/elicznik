package api

import (
	"elicznik/util"
	"encoding/csv"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

type ELicznikMeasurement struct {
	Date  string  `json:"Date"`
	Hour  int     `json:"Hour,string"`
	Power float64 `json:"EC,string"`
}

type ELicznikData struct {
	Ok   int `json:"ok"`
	Dane struct {
		FromGrid []ELicznikMeasurement `json:"chart"`
		FeedIn   []ELicznikMeasurement `json:"OZE"`
	} `json:"dane"`
}

func parseTauronCsvRecord(record []string) (*ELicznikMeasurement, error) {
	ts := strings.Split(record[0], " ")

	hour, err := strconv.Atoi(strings.Split(ts[1], ":")[0])
	if err != nil {
		return nil, err
	}

	date, err := time.Parse("2006-01-02", ts[0])
	if err != nil {
		return nil, err
	}

	power, err := strconv.ParseFloat(strings.ReplaceAll(record[1], ",", "."), 64)
	if err != nil {
		return nil, err
	}

	return &ELicznikMeasurement{date.Format(util.DATE_FORMAT), hour, power}, err
}

func ParseTauronCsv(payload string) (*ELicznikData, error) {
	r := csv.NewReader(strings.NewReader(payload))
	r.InputOffset()
	r.Comma = ';'
	r.Comment = '#'

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	// first record is csv header
	if len(records) > 1 {
		records = records[1:]
	}

	data := &ELicznikData{
		Ok: 1,
		Dane: struct {
			FromGrid []ELicznikMeasurement `json:"chart"`
			FeedIn   []ELicznikMeasurement `json:"OZE"`
		}{
			FromGrid: make([]ELicznikMeasurement, 0, len(records)),
			FeedIn:   make([]ELicznikMeasurement, 0, len(records)),
		},
	}

	for _, record := range records {
		eLicznikMeasurement, err := parseTauronCsvRecord(record)
		if err != nil {
			return nil, err
		}
		switch rt := record[2]; rt {
		case "oddanie":
			data.Dane.FeedIn = append(data.Dane.FeedIn, *eLicznikMeasurement)
		case "pobór":
			data.Dane.FromGrid = append(data.Dane.FromGrid, *eLicznikMeasurement)
		default:
			return nil, errors.Errorf("unknown csv record type '%v'", rt)
		}
	}

	return data, nil
}
