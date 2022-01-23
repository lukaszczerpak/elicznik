package api

import (
	"elicznik/util"
	"github.com/go-resty/resty/v2"
	"time"
)

type TauronApiClient struct {
	client                      *resty.Client
	username, password, smartNr string
}

func createRestyClient(apiKey string) *resty.Client {
	client := resty.New()
	//client.SetDebug(true)
	client.SetTimeout(1 * time.Minute)
	client.SetHeaders(map[string]string{
		"User-Agent": util.USER_AGENT,
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
