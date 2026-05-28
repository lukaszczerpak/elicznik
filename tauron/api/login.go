package api

import (
	"html"
	"net/url"
	"regexp"

	"github.com/pkg/errors"
)

var kcFormActionRe = regexp.MustCompile(`(?i)<form[^>]+(?:id="kc-form-login"|method="post")[^>]*\baction="([^"]+)"`)

func extractFormAction(body string) string {
	if m := kcFormActionRe.FindStringSubmatch(body); m != nil {
		return html.UnescapeString(m[1])
	}
	// fallback: any form with an action
	re := regexp.MustCompile(`(?i)<form[^>]+\baction="([^"]+)"`)
	if m := re.FindStringSubmatch(body); m != nil {
		return html.UnescapeString(m[1])
	}
	return ""
}

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

	// GET login page to obtain the dynamic Keycloak form action URL
	resp, err := api.client.R().Get("https://logowanie.tauron-dystrybucja.pl/login?service=https://elicznik.tauron-dystrybucja.pl")
	if err != nil {
		return errors.Wrap(err, "Session initialization error")
	}
	actionURL := extractFormAction(resp.String())
	if actionURL == "" {
		return errors.New("Could not extract login form action URL")
	}

	// POST credentials to the dynamic action URL
	_, err = api.client.R().
		SetFormData(map[string]string{
			"username":     api.username,
			"password":     api.password,
			"credentialId": "",
			"login":        "",
		}).
		Post(actionURL)
	if err != nil {
		return errors.Wrap(err, "Login error")
	}
	if !api.isUserLoggedIn() {
		return errors.New("Invalid username or password")
	}

	return nil
}
