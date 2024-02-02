package module

import (
	"context"
	"fmt"
	"strings"

	"github.com/fabled-se/zitadel-bootstraper/internal/bootstrap"
	"github.com/fabled-se/zitadel-bootstraper/internal/config"
	"github.com/fabled-se/zitadel-bootstraper/internal/zitadel"
	"github.com/rs/zerolog"
)

func NewAdminAccount(zClient *zitadel.Client, conf config.Config) bootstrap.Module {
	return &adminAccountModule{
		zClient:        zClient,
		zitadelOrgName: conf.Zitadel.OrgName,
		conf:           conf.AdminAccount,
	}
}

type adminAccountModule struct {
	zClient        *zitadel.Client
	zitadelOrgName string
	conf           config.AdminAccount
}

func (a *adminAccountModule) Name() string {
	return "AdminAccount"
}

func (a *adminAccountModule) Execute(ctx context.Context) error {
	log := zerolog.Ctx(ctx)

	if !a.conf.Setup {
		log.Warn().Msg("Setup is disabled, skipping module")
		return nil
	}

	org, err := a.zClient.GetOrgByName(a.zitadelOrgName)
	if err != nil {
		return fmt.Errorf("failed to get org by name: %w", err)
	}

	input := zitadel.CreateUserInput{
		OrgId:           org.Id,
		Username:        a.conf.Username,
		Firstname:       a.conf.Username,
		Lastname:        a.conf.Username,
		Email:           a.conf.Username,
		EmailIsVerified: true,
		Password:        a.conf.Password,
	}

	res, err := a.zClient.CreateUser(input)
	if err != nil {
		if strings.Contains(err.Error(), "409") {
			log.Warn().Str("username", a.conf.Username).Msg("User already exists")
			return nil
		}

		return fmt.Errorf("failed to create user: %w", err)
	}

	roles := []zitadel.IAMRole{
		zitadel.IAM_OWNER,
	}

	if err := a.zClient.AddIAMMember(res.UserId, roles); err != nil {
		return fmt.Errorf("failed to add IAM role to user: %w", err)
	}

	log.Info().
		Str("username", a.conf.Username).
		Str("userId", res.UserId).
		Msg("Successfully created admin account")

	return nil
}
