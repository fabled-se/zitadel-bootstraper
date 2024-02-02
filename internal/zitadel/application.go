package zitadel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateApplicationOIDCInput struct {
	OrgId                  string
	ProjectId              string
	Name                   string
	RedirectUris           []string
	ResponseTypes          []string
	GrantTypes             []string
	AppType                string
	AuthMethodType         string
	PostLogoutRedirectUris []string
	DevMode                bool
	AccessTokenType        string

	IdTokenRoleAssertion     bool
	IdTokenUserInfoAssertion bool
}

type CreateApplicationOIDCOutput struct {
	AppId        string `json:"appId"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func (c *Client) CreateApplicationOIDC(i CreateApplicationOIDCInput) (*CreateApplicationOIDCOutput, error) {
	payload := map[string]any{
		"name":                     i.Name,
		"redirectUris":             i.RedirectUris,
		"responseTypes":            i.ResponseTypes,
		"grantTypes":               i.GrantTypes,
		"appType":                  i.AppType,
		"authMethodType":           i.AuthMethodType,
		"postLogoutRedirectUris":   i.PostLogoutRedirectUris,
		"devMode":                  i.DevMode,
		"accessTokenType":          i.AccessTokenType,
		"idTokenRoleAssertion":     i.IdTokenRoleAssertion,
		"idTokenUserinfoAssertion": i.IdTokenUserInfoAssertion,
	}

	buffer := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buffer).Encode(payload); err != nil {
		return nil, fmt.Errorf("failed to encode payload: %w", err)
	}

	endpoint := fmt.Sprintf("management/v1/projects/%s/apps/oidc", i.ProjectId)
	req, err := c.newRequestWithAuth("POST", endpoint, buffer)
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

	var responseBody CreateApplicationOIDCOutput

	if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w")
	}

	return &responseBody, nil
}
