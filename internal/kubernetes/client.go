package kubernetes

import (
	"fmt"
	"io"
	"net/http"
)

func New(httpClient *http.Client, host, port string) *Client {
	return &Client{
		httpClient: httpClient,
		host:       host,
		port:       port,
	}
}

type Client struct {
	httpClient *http.Client
	host       string
	port       string

	namespace string
	token     string
}

func (c *Client) WithNamespace(namespace string) *Client {
	c.namespace = namespace
	return c
}

func (c *Client) WithToken(token string) *Client {
	c.token = token
	return c
}

func (c *Client) getBaseUrl() string {
	return fmt.Sprintf("https://%s:%s", c.host, c.port)
}

func (c *Client) newRequestWithAuth(method, endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.getBaseUrl()+"/"+endpoint, body)
	if err != nil {
		return req, err
	}

	req.Header.Add("authorization", "Bearer "+c.token)

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
