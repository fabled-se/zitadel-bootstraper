package zitadel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type GetOrgByNameOutput struct {
	Id   string
	Name string
}

func (c *Client) GetOrgByName(orgName string) (*GetOrgByNameOutput, error) {

	input := SearchOrgInput{
		OrgName: orgName,
		Offset:  0,
		Limit:   1,
		Asc:     true,
	}

	result, err := c.SearchOrg(input)
	if err != nil {
		return nil, err
	}

	if len(result.Orgs) == 0 {
		return nil, fmt.Errorf("no such org")
	}

	org := result.Orgs[0]

	// sanity
	if org.Name != orgName {
		return nil, fmt.Errorf("no such org")
	}

	return &GetOrgByNameOutput{Id: org.Id, Name: org.Name}, nil
}

type SearchOrgInput struct {
	OrgName string

	Offset int
	Limit  int
	Asc    bool
}

type SearchOrgOutput struct {
	Orgs []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"result"`
}

// https://zitadel.com/docs/apis/resources/admin/admin-service-list-orgs
func (c *Client) SearchOrg(i SearchOrgInput) (*SearchOrgOutput, error) {

	limit := 1000
	if i.Limit != 0 {
		limit = i.Limit
	}

	payload := map[string]any{
		"query": map[string]any{
			"offset": fmt.Sprintf("%d", i.Offset),
			"limit":  limit,
			"asc":    i.Asc,
		},
		"sortingColumn": "ORG_FIELD_NAME_UNSPECIFIED",
		"queries": []map[string]any{
			{
				"nameQuery": map[string]any{
					"name":   i.OrgName,
					"method": "TEXT_QUERY_METHOD_EQUALS",
				},
			},
		},
	}

	buffer := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buffer).Encode(payload); err != nil {
		return nil, fmt.Errorf("failed to encode payload: %w", err)
	}

	req, err := c.newRequestWithAuth("POST", "admin/v1/orgs/_search", buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if res.StatusCode != http.StatusOK {
		return nil, c.unexpectedStatusCodeErr(res.StatusCode, res.Body)
	}

	var responseBody SearchOrgOutput

	if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w")
	}

	return &responseBody, nil
}
