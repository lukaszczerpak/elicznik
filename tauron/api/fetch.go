package api

import (
	"elicznik/util"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

func (api *TauronApiClient) fetchData(dateFrom, dateTo time.Time) (*ELicznikData, error) {
	err := api.login()
	if err != nil {
		return nil, errors.Wrap(err, "Login error")
	}

	resp, err := api.client.R().
		SetFormData(map[string]string{
			"dane[chartDay]":  dateFrom.Format(util.DATE_FORMAT),
			"dane[startDay]":  dateFrom.Format(util.DATE_FORMAT),
			"dane[endDay]":    dateTo.Format(util.DATE_FORMAT),
			"dane[trybCSV]":   "godzin",
			"dane[paramType]": "csv",
			"dane[smartNr]":   api.smartNr,
			"dane[checkOZE]":  "on",
		}).
		//SetResult(&ELicznikData{}).
		Post("https://elicznik.tauron-dystrybucja.pl/index/charts")

	if err != nil {
		return nil, errors.Wrap(err, "Fetching data error")
	}

	// not possible since tauron serves json with Content-Type text/html...
	//result := resp.Result().(*ELicznikData)
	var result *ELicznikData
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, errors.Wrap(err, "JSON parsing error")
	}

	if result.Ok != 1 || len(result.Dane.FeedIn) != len(result.Dane.FromGrid) {
		return nil, errors.New(fmt.Sprintf("Data incomplete: ok=%v, len(FeedIn)=%v, len(FromGrid)=%v", result.Ok, len(result.Dane.FeedIn), len(result.Dane.FromGrid)))
	}

	return result, nil
}
