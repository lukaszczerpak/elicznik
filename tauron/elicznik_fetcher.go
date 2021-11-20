package tauron

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ELicznikFetcher struct {
	jar                         *cookiejar.Jar
	client                      *http.Client
	username, password, smartNr string
}

func NewELicznikFetcher(username, password, smartNr string) *ELicznikFetcher {
	fetcher := ELicznikFetcher{}
	fetcher.jar, _ = cookiejar.New(nil)
	fetcher.client = &http.Client{
		Jar: fetcher.jar,
	}
	fetcher.username = username
	fetcher.password = password
	fetcher.smartNr = smartNr
	return &fetcher
}

func (fetcher *ELicznikFetcher) isUserLoggedIn() bool {
	u, _ := url.Parse("https://elicznik.tauron-dystrybucja.pl")
	for _, cookie := range fetcher.jar.Cookies(u) {
		if cookie.Name == "PHPSESSID" {
			return true
		}
	}
	return false
}

func (fetcher *ELicznikFetcher) login() error {
	if fetcher.isUserLoggedIn() {
		return nil
	}

	req, _ := http.NewRequest("GET", "https://logowanie.tauron-dystrybucja.pl/login", nil)
	req.Header.Set("User-Agent", USER_AGENT)
	resp, err := fetcher.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Session initialization error")
	}
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	data := url.Values{}
	data.Set("username", fetcher.username)
	data.Set("password", fetcher.password)
	data.Set("service", "https://elicznik.tauron-dystrybucja.pl")
	req, _ = http.NewRequest("POST", "https://logowanie.tauron-dystrybucja.pl/login", strings.NewReader(data.Encode()))
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
	resp, err = fetcher.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Login error")
	}
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if !fetcher.isUserLoggedIn() {
		return errors.New("Invalid username or password")
	}
	return nil
}

func (fetcher *ELicznikFetcher) fetchData(dateFrom, dateTo time.Time) (*ELicznikData, error) {
	err := fetcher.login()
	if err != nil {
		return nil, errors.Wrap(err, "Login error")
	}

	data := url.Values{}
	data.Set("dane[chartDay]", dateFrom.Format(ELICZNIK_DATE_FORMAT))
	data.Set("dane[startDay]", dateFrom.Format(ELICZNIK_DATE_FORMAT))
	data.Set("dane[endDay]", dateTo.Format(ELICZNIK_DATE_FORMAT))
	data.Set("dane[trybCSV]", "godzin")
	data.Set("dane[paramType]", "csv")
	data.Set("dane[smartNr]", fetcher.smartNr)
	data.Set("dane[checkOZE]", "on")
	req, _ := http.NewRequest("POST", "https://elicznik.tauron-dystrybucja.pl/index/charts", strings.NewReader(data.Encode()))
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
	resp, err := fetcher.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Fetching data error")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// Unmarshal or Decode the JSON to the interface.
	var result ELicznikData
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	} else if result.Ok != 1 || len(result.Dane.PowerIn) != len(result.Dane.PowerOut) {
		return nil, errors.New(fmt.Sprintf("Data incomplete: ok=%s, len(powerIn)=%d, len(powerOut)=%d", result.Ok, len(result.Dane.PowerIn), len(result.Dane.PowerOut)))
	}

	return &result, nil
}
