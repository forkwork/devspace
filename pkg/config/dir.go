package config

import (
	"os"
	"path/filepath"

	"dev.khulnasoft.com/pkg/util"
)

// Override devspace home
const DEVPOD_HOME = "DEVPOD_HOME"

// Override config path
const DEVPOD_CONFIG = "DEVPOD_CONFIG"

func GetConfigDir() (string, error) {
	homeDir := os.Getenv(DEVPOD_HOME)
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
	configOrigin := os.Getenv(DEVPOD_CONFIG)
	if configOrigin == "" {
		configDir, err := GetConfigDir()
		if err != nil {
			return "", err
		}

		return filepath.Join(configDir, ConfigFile), nil
	}

	return configOrigin, nil
}
