package machine

import (
	"context"

	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/pkg/client"
	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/workspace"
	"dev.khulnasoft.com/log"
	"github.com/spf13/cobra"
)

// StopCmd holds the configuration
type StopCmd struct {
	*flags.GlobalFlags
}

// NewStopCmd creates a new destroy command
func NewStopCmd(flags *flags.GlobalFlags) *cobra.Command {
	cmd := &StopCmd{
		GlobalFlags: flags,
	}
	stopCmd := &cobra.Command{
		Use:   "stop [name]",
		Short: "Stops an existing machine",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.Run(context.Background(), args)
		},
	}

	return stopCmd
}

// Run runs the command logic
func (cmd *StopCmd) Run(ctx context.Context, args []string) error {
	devSpaceConfig, err := config.LoadConfig(cmd.Context, cmd.Provider)
	if err != nil {
		return err
	}

	machineClient, err := workspace.GetMachine(devSpaceConfig, args, log.Default)
	if err != nil {
		return err
	}

	err = machineClient.Stop(ctx, client.StopOptions{})
	if err != nil {
		return err
	}

	return nil
}
