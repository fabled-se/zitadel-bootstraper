package zitadel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateUserInput struct {
	OrgId           string
	Username        string
	Firstname       string
	Lastname        string
	Email           string
	EmailIsVerified bool
	Password        string
}

type CreateUserOutput struct {
	UserId string `json:"userId"`
}

// https://zitadel.com/docs/apis/resources/mgmt/management-service-import-human-user
func (c *Client) CreateUser(i CreateUserInput) (*CreateUserOutput, error) {
	payload := map[string]any{
		"userName": i.Username,
		"profile": map[string]any{
			"firstName": i.Firstname,
			"lastName":  i.Lastname,
		},
		"email": map[string]any{
			"email":           i.Email,
			"isEmailVerified": i.EmailIsVerified,
		},
		"password": i.Password,
	}

	buffer := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buffer).Encode(payload); err != nil {
		return nil, fmt.Errorf("failed to encode payload: %w", err)
	}

	req, err := c.newRequestWithAuth("POST", "management/v1/users/human/_import", buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("x-zitadel-orgid", i.OrgId)

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, c.unexpectedStatusCodeErr(res)
	}

	var responseBody CreateUserOutput

	if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &responseBody, nil
}
