package platform

import (
	"context"
	"fmt"
	"path/filepath"

	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/platform/client"
	"dev.khulnasoft.com/pkg/provider"
	"dev.khulnasoft.com/log"
)

func InitClientFromHost(ctx context.Context, devSpaceConfig *config.Config, devSpaceProHost string, log log.Logger) (client.Client, error) {
	provider, err := ProviderFromHost(ctx, devSpaceConfig, devSpaceProHost, log)
	if err != nil {
		return nil, fmt.Errorf("provider from pro instance: %w", err)
	}

	return InitClientFromProvider(ctx, devSpaceConfig, provider, log)
}

func InitClientFromProvider(ctx context.Context, devSpaceConfig *config.Config, providerName string, log log.Logger) (client.Client, error) {
	configPath, err := LoftConfigPath(devSpaceConfig.DefaultContext, providerName)
	if err != nil {
		return nil, fmt.Errorf("loft config path: %w", err)
	}

	return client.InitClientFromPath(ctx, configPath)
}

func ProviderFromHost(ctx context.Context, devSpaceConfig *config.Config, devSpaceProHost string, log log.Logger) (string, error) {
	proInstanceConfig, err := provider.LoadProInstanceConfig(devSpaceConfig.DefaultContext, devSpaceProHost)
	if err != nil {
		return "", fmt.Errorf("load pro instance %s: %w", devSpaceProHost, err)
	}

	return proInstanceConfig.Provider, nil
}

func LoftConfigPath(context string, providerName string) (string, error) {
	providerDir, err := provider.GetProviderDir(context, providerName)
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(providerDir, "loft-config.json")

	return configPath, nil
}
