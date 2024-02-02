package zitadel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/fabled-se/zitadel-bootstraper/internal/config"
)

func New(httpClient *http.Client, conf config.Zitadel, keyJson string) (*Client, error) {
	c := &Client{
		HttpClient:  httpClient,
		TLS:         conf.TLS,
		Domain:      conf.Domain,
		OrgName:     conf.OrgName,
		ServiceUser: conf.ServiceUserName,
	}

	jwt, err := c.newJWT([]byte(keyJson), conf.Domain)
	if err != nil {
		return nil, fmt.Errorf("failed to create jwt: %w", err)
	}

	if err := c.requestOauthToken(jwt); err != nil {
		return nil, fmt.Errorf("failed to setup oauth token: %w", err)
	}

	return c, nil
}

type Client struct {
	HttpClient  *http.Client
	TLS         bool
	Domain      string
	OrgName     string
	ServiceUser string

	serviceUserToken string
}

func (c *Client) newRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, c.getBaseUrl()+"/"+endpoint, body)
}

func (c *Client) newRequestWithAuth(method, endpoint string, body io.Reader) (*http.Request, error) {
	req, err := c.newRequest(method, endpoint, body)
	if err != nil {
		return req, err
	}

	req.Header.Add("Authorization", "Bearer "+c.serviceUserToken)

	return req, nil
}

func (c *Client) unexpectedStatusCodeErr(res *http.Response) error {
	bodyBytes, _ := io.ReadAll(res.Body)
	return fmt.Errorf(
		"unexpected status code %d, response body: %s",
		res.StatusCode,
		string(bodyBytes),
	)
}

func (c *Client) getBaseUrl() string {
	protocol := "http"
	if c.TLS {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s", protocol, c.Domain)
}

func (c *Client) requestOauthToken(jwt string) error {
	form := url.Values{}
	form.Add("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	form.Add("scope", "openid profile email urn:zitadel:iam:org:project:id:zitadel:aud")
	form.Add("assertion", jwt)

	req, err := c.newRequest("POST", "oauth/v2/token", strings.NewReader(form.Encode()))
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
		return c.unexpectedStatusCodeErr(res)
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
