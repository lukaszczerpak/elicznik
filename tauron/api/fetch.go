package api

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

const DATE_FORMAT = "02.01.2006"

func (api *TauronApiClient) fetchData(dateFrom, dateTo time.Time) (*ELicznikData, error) {
	err := api.login()
	if err != nil {
		return nil, errors.Wrap(err, "Login error")
	}

	resp, err := api.client.R().
		SetQueryParams(map[string]string{
			"form[from]":     dateFrom.Format(DATE_FORMAT),
			"form[to]":       dateTo.Format(DATE_FORMAT),
			"form[type]":     "godzin",
			"form[consum]":   "1",
			"form[oze]":      "1",
			"form[fileType]": "CSV",
		}).
		Get("https://elicznik.tauron-dystrybucja.pl/energia/do/dane")

	if err != nil {
		return nil, errors.Wrap(err, "Fetching data error")
	}

	var result *ELicznikData
	result, err = ParseTauronCsv(resp.String())
	if err != nil {
		return nil, errors.Wrap(err, "CSV parsing error")
	}

	if result.Ok != 1 || len(result.Dane.FeedIn) != len(result.Dane.FromGrid) {
		return nil, errors.New(fmt.Sprintf("Data incomplete: ok=%v, len(FeedIn)=%v, len(FromGrid)=%v", result.Ok, len(result.Dane.FeedIn), len(result.Dane.FromGrid)))
	}

	return result, nil
}
