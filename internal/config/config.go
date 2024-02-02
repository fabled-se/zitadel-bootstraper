package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Zitadel Zitadel `yaml:"zitadel"`

	AdminAccount AdminAccount `yaml:"adminAccount"`
}

type Zitadel struct {
	Domain          string `yaml:"domain"`
	TLS             bool   `yaml:"tls"`
	OrgName         string `yaml:"orgName"`
	ServiceUserName string `yaml:"serviceUserName"`
}

type AdminAccount struct {
	Setup    bool   `yaml:"setup"`
	OrgName  string `yaml:"orgName"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func ParseFromFile(path string) (Config, error) {
	var c Config

	f, err := os.Open(path)
	if err != nil {
		return c, fmt.Errorf("failed to open file path: %w", err)
	}

	if err := yaml.NewDecoder(f).Decode(&c); err != nil {
		return c, fmt.Errorf("failed to decode yaml: %w", err)
	}

	return c, nil
}
