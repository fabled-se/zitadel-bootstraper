package zitadel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateProjectInput struct {
	OrgId                string
	Name                 string
	ProjectRoleAssertion bool
	ProjectRoleCheck     bool
	HasProjectCheck      bool
}

type CreateProjectOutput struct {
	Id string `json:"id"`
}

func (c *Client) CreateProject(i CreateProjectInput) (*CreateProjectOutput, error) {
	payload := map[string]any{
		"name":                 i.Name,
		"projectRoleAssertion": i.ProjectRoleAssertion,
		"projectRoleCheck":     i.ProjectRoleCheck,
		"hasProjectCheck":      i.HasProjectCheck,
	}

	buffer := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buffer).Encode(payload); err != nil {
		return nil, fmt.Errorf("failed to encode payload: %w", err)
	}

	req, err := c.newRequestWithAuth("POST", "management/v1/projects", buffer)
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

	var responseBody CreateProjectOutput

	if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &responseBody, nil
}

type BulkAddProjectRoleInput struct {
	OrgId     string
	ProjectId string
	Roles     []ProjectRole
}

type ProjectRole struct {
	Key         string `json:"key"`
	DisplayName string `json:"displayName"`
	Group       string `json:"group"`
}

// https://zitadel.com/docs/apis/resources/mgmt/management-service-bulk-add-project-roles
func (c *Client) BulkAddProjectRole(i BulkAddProjectRoleInput) error {
	payload := map[string]any{
		"roles": i.Roles,
	}

	buffer := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buffer).Encode(payload); err != nil {
		return fmt.Errorf("failed to encode payload: %w", err)
	}

	endpoint := fmt.Sprintf("management/v1/projects/%s/roles/_bulk", i.ProjectId)
	req, err := c.newRequestWithAuth("POST", endpoint, buffer)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("x-zitadel-orgid", i.OrgId)

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return c.unexpectedStatusCodeErr(res)
	}

	return nil
}
