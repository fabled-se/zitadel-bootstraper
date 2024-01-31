package zitadel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	HttpClient  *http.Client
	TLS         bool
	Domain      string
	OrgName     string
	ServiceUser string

	serviceUserToken string
}

func (c *Client) getBaseUrl() string {
	protocol := "http"
	if c.TLS {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s", protocol, c.Domain)
}

func (c *Client) SetupOauthToken(jwt string) error {
	form := url.Values{}
	form.Add("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	form.Add("scope", "openid profile email")
	form.Add("assertion", jwt)

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/oauth/v2/token", c.getBaseUrl()),
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perfrom request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("unexpected status code %d, body: %s", res.StatusCode, string(bodyBytes))
	}

	tokenResponse := struct {
		AccessToken string `json:"access_token"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&tokenResponse); err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}

	c.serviceUserToken = tokenResponse.AccessToken

	return nil
}
