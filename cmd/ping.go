package cmd

import (
	"context"
	"fmt"
	"os"

	"dev.khulnasoft.com/cmd/completion"
	"dev.khulnasoft.com/cmd/flags"
	client2 "dev.khulnasoft.com/pkg/client"
	"dev.khulnasoft.com/pkg/config"
	workspace2 "dev.khulnasoft.com/pkg/workspace"
	"dev.khulnasoft.com/log"
	"github.com/spf13/cobra"
)

type PingCmd struct {
	*flags.GlobalFlags
}

func NewPingCmd(flags *flags.GlobalFlags) *cobra.Command {
	cmd := &PingCmd{
		GlobalFlags: flags,
	}
	troubleshootCmd := &cobra.Command{
		Use:   "ping [workspace-path|workspace-name]",
		Short: "Pings the DevSpace Pro workspace",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.Run(cobraCmd.Context(), args)
		},
		ValidArgsFunction: func(rootCmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completion.GetWorkspaceSuggestions(rootCmd, cmd.Context, cmd.Provider, args, toComplete, cmd.Owner, log.Default)
		},
		Hidden: true,
	}

	return troubleshootCmd
}

func (cmd *PingCmd) Run(ctx context.Context, args []string) error {
	devSpaceConfig, err := config.LoadConfig(cmd.Context, cmd.Provider)
	if err != nil {
		return err
	}

	client, err := workspace2.Get(ctx, devSpaceConfig, args, true, cmd.Owner, false, log.Default.ErrorStreamOnly())
	if err != nil {
		return err
	}

	daemonClient, ok := client.(client2.DaemonClient)
	if !ok {
		return fmt.Errorf("ping is only available for pro workspaces")
	}

	return daemonClient.Ping(ctx, os.Stdout)
}
