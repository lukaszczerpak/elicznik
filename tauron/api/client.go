package api

import (
	"github.com/go-resty/resty/v2"
	"time"
)

const UA = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
const REFERER = "https://elicznik.tauron-dystrybucja.pl/energia"

type TauronApiClient struct {
	client                      *resty.Client
	username, password, smartNr string
}

func createRestyClient(apiKey string) *resty.Client {
	client := resty.New()
	//client.SetDebug(true)
	client.SetTimeout(1 * time.Minute)
	client.SetHeaders(map[string]string{
		"User-Agent": UA,
		"Referer":    REFERER,
	})

	return client
}

func New(username, password, smartNr string) *TauronApiClient {
	apiClient := TauronApiClient{}
	apiClient.username = username
	apiClient.password = password
	apiClient.smartNr = smartNr
	apiClient.client = resty.New()
	return &apiClient
}
