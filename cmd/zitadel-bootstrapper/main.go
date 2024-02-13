package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"os"

	"github.com/fabled-se/zitadel-bootstraper/internal/bootstrap"
	"github.com/fabled-se/zitadel-bootstraper/internal/config"
	"github.com/fabled-se/zitadel-bootstraper/internal/kubernetes"
	"github.com/fabled-se/zitadel-bootstraper/internal/module"
	"github.com/fabled-se/zitadel-bootstraper/internal/zitadel"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout)

	kubernetesApiHost := mustEnvVar(logger, "KUBERNETES_SERVICE_HOST")
	kubernetesAPiPort := mustEnvVar(logger, "KUBERNETES_PORT_443_TCP_PORT")
	kubernetesToken := mustFileReadAll(logger, "var/run/secrets/kubernetes.io/serviceaccount/token")
	kubernetsCACert := mustFileReadAll(logger, "var/run/secrets/kubernetes.io/serviceaccount/ca.crt")

	zitadelServiceUserKeyJson := mustEnvVar(logger, "ZITADEL_SERVICE_USER_KEY_JSON")

	bootstrapConfig, err := config.ParseFromFile("/etc/zitadel-bootstrapper-config/config-yaml")
	if err != nil {
		logger.Err(err).Msg("Failed to parse config file")
		os.Exit(1)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(kubernetsCACert)

	kHttpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	// TODO: take namespace as helm value?
	kClient := kubernetes.New(kHttpClient, kubernetesApiHost, kubernetesAPiPort).
		WithNamespace("argocd").
		WithToken(string(kubernetesToken))

	zClient, err := zitadel.New(http.DefaultClient, bootstrapConfig.Zitadel, zitadelServiceUserKeyJson)
	if err != nil {
		logger.Err(err).Msg("Failed to create zitadel client")
		os.Exit(1)
	}

	modules := []bootstrap.Module{
		module.NewAdminAccount(zClient, bootstrapConfig),
		module.NewArgoCD(zClient, kClient, bootstrapConfig),
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

func mustFileReadAll(logger zerolog.Logger, path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		logger.Err(err).Str("path", path).Msg("Failed to open file")
		os.Exit(1)
	}

	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		logger.Err(err).Msg("Failed to read bytes from file")
		os.Exit(1)
	}

	return b
}

func mustEnvVar(logger zerolog.Logger, key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Error().Msgf("Environment variable %s is empty!", key)
		os.Exit(1)
	}

	return value
}
