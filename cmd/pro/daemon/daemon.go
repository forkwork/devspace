package daemon

import (
	"context"
	"fmt"

	"dev.khulnasoft.com/cmd/pro/flags"
	"dev.khulnasoft.com/pkg/config"
	providerpkg "dev.khulnasoft.com/pkg/provider"
	"dev.khulnasoft.com/pkg/workspace"
	"dev.khulnasoft.com/log"
	"github.com/spf13/cobra"
)

// NewCmd creates a new cobra command
func NewCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	c := &cobra.Command{
		Use:    "daemon",
		Short:  "DevSpace Pro Provider daemon commands",
		Args:   cobra.NoArgs,
		Hidden: true,
	}

	c.AddCommand(NewStartCmd(globalFlags))
	c.AddCommand(NewStatusCmd(globalFlags))
	c.AddCommand(NewNetcheckCmd(globalFlags))

	return c
}

func findProProvider(ctx context.Context, context, provider, host string, log log.Logger) (*config.Config, *providerpkg.ProviderConfig, error) {
	devSpaceConfig, err := config.LoadConfig(context, provider)
	if err != nil {
		return nil, nil, err
	}

	pCfg, err := workspace.ProviderFromHost(ctx, devSpaceConfig, host, log)
	if err != nil {
		return devSpaceConfig, nil, fmt.Errorf("load provider: %w", err)
	}

	return devSpaceConfig, pCfg, nil
}
