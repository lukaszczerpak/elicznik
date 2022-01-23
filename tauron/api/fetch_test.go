package api

import (
	"encoding/json"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func addLoginCookie(c *http.Client) {
	u, _ := url.Parse("https://elicznik.tauron-dystrybucja.pl")
	cookies := []*http.Cookie{
		&http.Cookie{
			Name:  "PHPSESSID",
			Value: "This is cookie value",
		}}
	c.Jar.SetCookies(u, cookies)
}

func loadFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func loadJsonToStruct(filename string, v interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return json.Unmarshal(data, &v)
}

func mockJsonResponse(filename string) {
	httpmock.RegisterResponder("GET", `https://logowanie.tauron-dystrybucja.pl/login`,
		httpmock.NewJsonResponderOrPanic(200, make(map[string]string)))
	httpmock.RegisterResponder("POST", `https://logowanie.tauron-dystrybucja.pl/login`,
		httpmock.NewJsonResponderOrPanic(200, make(map[string]string)))
	httpmock.RegisterResponder("POST", `https://elicznik.tauron-dystrybucja.pl/index/charts`,
		httpmock.NewBytesResponder(200, loadFile(filename)))

	//httpmock.RegisterResponder("POST", `https://elicznik.tauron-dystrybucja.pl/index/charts`,
	//	httpmock.NewJsonResponderOrPanic(200, util.LoadJsonToMap(filename)))
}

func TestFetchData(t *testing.T) {
	client := createRestyClient("")
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	addLoginCookie(client.GetClient())
	mockJsonResponse("../testdata/sample.json")

	s := &TauronApiClient{client: client}
	resp, err := s.fetchData(time.Now(), time.Now())

	assert.NoError(t, err, "expected no error")
	assert.Equal(t, 1, httpmock.GetTotalCallCount(), "number of api calls")
	assert.Equal(t, 1, resp.Ok)
	assert.Len(t, resp.Dane.FeedIn, 168, "dane.feedin")
	assert.Len(t, resp.Dane.FromGrid, 168, "dane.fromgrid")
}

func TestFetchDataCET2CEST(t *testing.T) {
	client := createRestyClient("")
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	addLoginCookie(client.GetClient())
	mockJsonResponse("../testdata/sample-CET-to-CEST.json")

	s := &TauronApiClient{client: client}
	resp, err := s.fetchData(time.Now(), time.Now())

	assert.NoError(t, err, "expected no error")
	assert.Equal(t, 1, httpmock.GetTotalCallCount(), "number of api calls")
	assert.Equal(t, 1, resp.Ok)
	assert.Len(t, resp.Dane.FeedIn, 23, "dane.feedin")
	assert.Len(t, resp.Dane.FromGrid, 23, "dane.fromgrid")
}

func TestFetchDataCEST2CET(t *testing.T) {
	client := createRestyClient("")
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	addLoginCookie(client.GetClient())
	mockJsonResponse("../testdata/sample-CEST-to-CET.json")

	s := &TauronApiClient{client: client}
	resp, err := s.fetchData(time.Now(), time.Now())

	assert.NoError(t, err, "expected no error")
	assert.Equal(t, 1, httpmock.GetTotalCallCount(), "number of api calls")
	assert.Equal(t, 1, resp.Ok)
	assert.Len(t, resp.Dane.FeedIn, 25, "dane.feedin")
	assert.Len(t, resp.Dane.FromGrid, 25, "dane.fromgrid")
}
