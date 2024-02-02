package main

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/fabled-se/zitadel-bootstraper/internal/bootstrap"
	"github.com/fabled-se/zitadel-bootstraper/internal/config"
	"github.com/fabled-se/zitadel-bootstraper/internal/module"
	"github.com/fabled-se/zitadel-bootstraper/internal/zitadel"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout)

	zitadelServiceUserKeyJson := mustEnvVar(logger, "ZITADEL_SERVICE_USER_KEY_JSON")

	bootstrapConfig, err := config.ParseFromFile("/etc/zitadel-bootstrapper-config/config-yaml")
	if err != nil {
		logger.Err(err).Msg("Failed to parse config file")
		os.Exit(1)
	}

	zClient, err := zitadel.New(http.DefaultClient, bootstrapConfig.Zitadel, zitadelServiceUserKeyJson)
	if err != nil {
		logger.Err(err).Msg("Failed to create zitadel client")
		os.Exit(1)
	}

	modules := []bootstrap.Module{
		module.NewAdminAccount(zClient, bootstrapConfig),
		module.NewArgoCD(zClient, bootstrapConfig),
	}

	for _, module := range modules {
		log := logger.With().Str("module", module.Name()).Logger()

		// TODO: Context with deadline?
		if err := module.Execute(log.WithContext(context.TODO())); err != nil {
			log.Err(err).Msg("Failed to execute module")
			os.Exit(1)
		}
	}

	logger.Info().Msg("Bootstrapping successful")
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
