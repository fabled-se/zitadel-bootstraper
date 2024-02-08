package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Zitadel Zitadel `yaml:"zitadel"`

	AdminAccount AdminAccount `yaml:"adminAccount"`
	ArgoCD       ArgoCD       `yaml:"argoCD"`
}

type Zitadel struct {
	Domain          string `yaml:"domain"`
	TLS             bool   `yaml:"tls"`
	OrgName         string `yaml:"orgName"`
	ServiceUserName string `yaml:"serviceUserName"`
}

type AdminAccount struct {
	Setup     bool   `yaml:"setup"`
	Firstname string `yaml:"firstname"`
	Lastname  string `yaml:"lastname"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

type ArgoCD struct {
	Setup         bool     `yaml:"setup"`
	Name          string   `yaml:"name"`
	UserRoleName  string   `yaml:"userRoleName"`
	AdminRoleName string   `yaml:"adminRoleName"`
	DevMode       bool     `yaml:"devMode"`
	RedirectUris  []string `yaml:"redirectUris"`
	LogoutUris    []string `yaml:"logoutUris"`
}

func ParseFromFile(path string) (Config, error) {
	var c Config

	f, err := os.Open(path)
	if err != nil {
		return c, fmt.Errorf("failed to open file path: %w", err)
	}

	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&c); err != nil {
		return c, fmt.Errorf("failed to decode yaml: %w", err)
	}

	return c, nil
}
