package zitadel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type IAMRole string

const (
	IAM_OWNER IAMRole = "IAM_OWNER"
)

func (c *Client) AddIAMMember(userId string, roles []IAMRole) error {
	payload := map[string]any{
		"userId": userId,
		"roles":  roles,
	}

	buffer := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buffer).Encode(payload); err != nil {
		return fmt.Errorf("failed to encode payload: %w", err)
	}

	req, err := c.newRequestWithAuth("POST", "admin/v1/members", buffer)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return c.unexpectedStatusCodeErr(res)
	}

	return nil
}
