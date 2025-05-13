package context

import (
	"context"
	"fmt"

	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/pkg/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// UseCmd holds the use cmd flags
type UseCmd struct {
	flags.GlobalFlags

	Options []string
}

// NewUseCmd uses a new command
func NewUseCmd(flags *flags.GlobalFlags) *cobra.Command {
	cmd := &UseCmd{
		GlobalFlags: *flags,
	}
	useCmd := &cobra.Command{
		Use:   "use",
		Short: "Set a DevSpace context as the default",
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("please specify the context to use")
			}

			return cmd.Run(context.Background(), args[0])
		},
	}

	useCmd.Flags().StringArrayVarP(&cmd.Options, "option", "o", []string{}, "context option in the form KEY=VALUE")
	return useCmd
}

// Run runs the command logic
func (cmd *UseCmd) Run(ctx context.Context, context string) error {
	devSpaceConfig, err := config.LoadConfig("", cmd.Provider)
	if err != nil {
		return err
	} else if devSpaceConfig.Contexts[context] == nil {
		return fmt.Errorf("context '%s' doesn't exist", context)
	}

	// check if there are use options set
	if len(cmd.Options) > 0 {
		err = setOptions(devSpaceConfig, context, cmd.Options)
		if err != nil {
			return err
		}
	}

	devSpaceConfig.DefaultContext = context
	err = config.SaveConfig(devSpaceConfig)
	if err != nil {
		return errors.Wrap(err, "save config")
	}

	return nil
}
