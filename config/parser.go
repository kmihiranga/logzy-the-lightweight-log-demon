package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Parse parses all configuration to a single config object
func Parse() *Config {
	return &Config{
		AppConfig: parseAppConfig(),
	}
}

func parseAppConfig() *AppConfig {
	configFile := CheckEnvironment("app")
	content := ReadFile(configFile)

	cfg := &AppConfig{}

	err := yaml.Unmarshal(*content, &cfg)
	if err != nil {
		panic(fmt.Sprintf("error: %v", err))
	}

	return cfg
}

// check development environment
func CheckEnvironment(folderName string) string {
	env := os.Getenv("ENV")
	var configFileName string

	switch env {
	case "production":
		configFileName = "/config-production.yaml"
	case "staging":
		configFileName = "/config-staging.yaml"
	case "development":
		configFileName = "/config-development.yaml"
	default:
		configFileName = "/config-local.yaml"
	}

	configFilePath := filepath.Join(folderName, configFileName)
	return configFilePath
}
