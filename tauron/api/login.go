package api

import (
	"github.com/pkg/errors"
	"net/url"
)

func (api *TauronApiClient) isUserLoggedIn() bool {
	u, _ := url.Parse("https://elicznik.tauron-dystrybucja.pl")
	for _, cookie := range api.client.GetClient().Jar.Cookies(u) {
		if cookie.Name == "PHPSESSID" {
			return true
		}
	}
	return false
}

func (api *TauronApiClient) login() error {
	if api.isUserLoggedIn() {
		return nil
	}

	// initialize session
	_, err := api.client.R().Get("https://logowanie.tauron-dystrybucja.pl/login")
	if err != nil {
		return errors.Wrap(err, "Session initialization error")
	}

	// send login form
	_, err = api.client.R().
		SetFormData(map[string]string{
			"username": api.username,
			"password": api.password,
			"service":  "https://elicznik.tauron-dystrybucja.pl",
		}).
		Post("https://logowanie.tauron-dystrybucja.pl/login")
	if err != nil {
		return errors.Wrap(err, "Login error")
	}
	if !api.isUserLoggedIn() {
		return errors.New("Invalid username or password")
	}

	// all good
	return nil
}
