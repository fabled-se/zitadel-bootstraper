package kubernetes

import (
	"bytes"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v3"
)

func (c *Client) CreateSecretStringData(name string, labels, values map[string]string) error {
	payload := map[string]any{
		"apiVersion": "v1",
		"kind":       "Secret",
		"metadata": map[string]any{
			"name":   name,
			"labels": labels,
		},
		"stringData": values,
	}

	buffer := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(buffer).Encode(payload); err != nil {
		return fmt.Errorf("failed to encode payload: %w", err)
	}

	endpoint := fmt.Sprintf("api/v1/namespaces/%s/secrets", c.namespace)
	req, err := c.newRequestWithAuth("POST", endpoint, buffer)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("content-type", "application/yaml")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return c.unexpectedStatusCodeErr(res)
	}

	return nil
}
