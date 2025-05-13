package config

import (
	"os"
	"path/filepath"

	"dev.khulnasoft.com/pkg/util"
)

// Override devspace home
const DEVSPACE_HOME = "DEVSPACE_HOME"

// Override config path
const DEVSPACE_CONFIG = "DEVSPACE_CONFIG"

func GetConfigDir() (string, error) {
	homeDir := os.Getenv(DEVSPACE_HOME)
	if homeDir != "" {
		return homeDir, nil
	}

	homeDir, err := util.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".devspace")
	return configDir, nil
}

func GetConfigPath() (string, error) {
	configOrigin := os.Getenv(DEVSPACE_CONFIG)
	if configOrigin == "" {
		configDir, err := GetConfigDir()
		if err != nil {
			return "", err
		}

		return filepath.Join(configDir, ConfigFile), nil
	}

	return configOrigin, nil
}
