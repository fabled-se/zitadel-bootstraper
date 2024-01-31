package main

import (
	"net/http"
	"os"

	"github.com/fabled-se/zitadel-bootstraper/internal/zitadel"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout)

	zitadelDomain := mustEnvVar(logger, "ZITADEL_DOMAIN")
	zitadelOrgName := mustEnvVar(logger, "ZITADEL_ORGNAME")
	zitadelServiceUser := mustEnvVar(logger, "ZITADEL_SERVICE_USER")
	zitadelServiceUserKeyJson := mustEnvVar(logger, "ZITADEL_SERVICE_USER_KEY_JSON")

	jwtKey, err := zitadel.NewJWT([]byte(zitadelServiceUserKeyJson), zitadelDomain)
	if err != nil {
		logger.Err(err).Msg("Failed to create a new Zitadel JWT")
		os.Exit(1)
	}

	zitadelClient := zitadel.Client{
		HttpClient:  http.DefaultClient,
		Domain:      zitadelDomain,
		OrgName:     zitadelOrgName,
		ServiceUser: zitadelServiceUser,
	}

	if err := zitadelClient.SetupOauthToken(jwtKey); err != nil {
		logger.Err(err).Msg("Failed to setup zitadel oauth token")
		os.Exit(1)
	}

	logger.Info().Msg("Success")
}

func mustEnvVar(logger zerolog.Logger, key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Error().Msgf("Environment variable %s is empty!", key)
		os.Exit(1)
	}

	return value
}
