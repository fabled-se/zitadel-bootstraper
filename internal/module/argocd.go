package module

import (
	"context"
	"fmt"

	"github.com/fabled-se/zitadel-bootstraper/internal/bootstrap"
	"github.com/fabled-se/zitadel-bootstraper/internal/config"
	"github.com/fabled-se/zitadel-bootstraper/internal/kubernetes"
	"github.com/fabled-se/zitadel-bootstraper/internal/zitadel"
	"github.com/rs/zerolog"
)

func NewArgoCD(zClient *zitadel.Client, kClient *kubernetes.Client, conf config.Config) bootstrap.Module {
	return &argocdModule{
		zClient:        zClient,
		kClient:        kClient,
		zitadelOrgName: conf.Zitadel.OrgName,
		conf:           conf.ArgoCD,
	}
}

type argocdModule struct {
	zClient        *zitadel.Client
	kClient        *kubernetes.Client
	zitadelOrgName string
	conf           config.ArgoCD
}

func (a *argocdModule) Name() string {
	return "ArgoCD"
}

func (a *argocdModule) Execute(ctx context.Context) error {
	log := zerolog.Ctx(ctx)

	if !a.conf.Setup {
		log.Warn().Msg("Setup is disabled, skipping module")
		return nil
	}

	org, err := a.zClient.GetOrgByName(a.zitadelOrgName)
	if err != nil {
		return fmt.Errorf("failed to get org by name: %w", err)
	}

	project, err := a.zClient.CreateProject(zitadel.CreateProjectInput{
		OrgId:                org.Id,
		Name:                 "ArgoCD",
		ProjectRoleAssertion: true,
		ProjectRoleCheck:     true,
		HasProjectCheck:      true,
	})
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	roles := []zitadel.ProjectRole{
		{Key: a.conf.UserRoleName, DisplayName: a.conf.UserRoleName},
		{Key: a.conf.AdminRoleName, DisplayName: a.conf.AdminRoleName},
	}

	bulkRoleInput := zitadel.BulkAddProjectRoleInput{
		OrgId:     org.Id,
		ProjectId: project.Id,
		Roles:     roles,
	}

	if err := a.zClient.BulkAddProjectRole(bulkRoleInput); err != nil {
		return fmt.Errorf("failed to add project roles: %w", err)
	}

	applicationInput := zitadel.CreateApplicationOIDCInput{
		OrgId:                    org.Id,
		ProjectId:                project.Id,
		Name:                     a.conf.Name,
		RedirectUris:             a.conf.RedirectUris,
		ResponseTypes:            []string{"OIDC_RESPONSE_TYPE_CODE"},
		GrantTypes:               []string{"OIDC_GRANT_TYPE_AUTHORIZATION_CODE"},
		AppType:                  "OIDC_APP_TYPE_WEB",
		AuthMethodType:           "OIDC_AUTH_METHOD_TYPE_BASIC",
		PostLogoutRedirectUris:   a.conf.LogoutUris,
		DevMode:                  a.conf.DevMode,
		AccessTokenType:          "OIDC_TOKEN_TYPE_BEARER",
		IdTokenRoleAssertion:     true,
		IdTokenUserInfoAssertion: true,
	}

	application, err := a.zClient.CreateApplicationOIDC(applicationInput)
	if err != nil {
		return fmt.Errorf("failed to create oidc application: %w", err)
	}

	labels := map[string]string{"app.kubernetes.io/part-of": "argocd"}
	values := map[string]string{
		"clientId": application.ClientId,
		"secretId": application.ClientSecret,
	}

	if err := a.kClient.CreateSecretStringData("argocd-zitadel-oidc", labels, values); err != nil {
		return fmt.Errorf("failed to create kubernetes oidc secret: %w", err)
	}

	log.Info().Interface("application", application).Msg("Created ArgoCD project")

	return nil
}
