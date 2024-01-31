package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/fabled-se/zitadel-bootstraper/internal/zitadel"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout)

	zitadelTLS := mustBool(logger, mustEnvVar(logger, "ZITADEL_TLS"))
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
		TLS:         zitadelTLS,
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

func mustBool(logger zerolog.Logger, value string) bool {
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		logger.Err(err).Msgf("Value '%s' must be parsed to bool", value)
		os.Exit(1)
	}
	return boolValue
}

func mustEnvVar(logger zerolog.Logger, key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Error().Msgf("Environment variable %s is empty!", key)
		os.Exit(1)
	}

	return value
}
