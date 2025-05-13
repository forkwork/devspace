package context

import (
	"context"
	"fmt"

	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/pkg/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// DeleteCmd holds the delete cmd flags
type DeleteCmd struct {
	flags.GlobalFlags
}

// NewDeleteCmd deletes a new command
func NewDeleteCmd(flags *flags.GlobalFlags) *cobra.Command {
	cmd := &DeleteCmd{
		GlobalFlags: *flags,
	}
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a DevSpace context",
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) > 1 {
				return fmt.Errorf("please specify the context to delete")
			}

			devSpaceContext := ""
			if len(args) == 1 {
				devSpaceContext = args[0]
			}

			return cmd.Run(context.Background(), devSpaceContext)
		},
	}

	return deleteCmd
}

// Run runs the command logic
func (cmd *DeleteCmd) Run(ctx context.Context, context string) error {
	devSpaceConfig, err := config.LoadConfig(context, cmd.Provider)
	if err != nil {
		return err
	}

	// check for context
	if context == "" {
		context = devSpaceConfig.DefaultContext
	} else if devSpaceConfig.Contexts[context] == nil {
		return fmt.Errorf("context '%s' doesn't exist", context)
	}

	// check for default context
	if context == "default" {
		return fmt.Errorf("cannot delete 'default' context")
	}

	delete(devSpaceConfig.Contexts, context)
	if devSpaceConfig.DefaultContext == context {
		devSpaceConfig.DefaultContext = "default"
	}
	if devSpaceConfig.OriginalContext == context {
		devSpaceConfig.OriginalContext = "default"
	}

	err = config.SaveConfig(devSpaceConfig)
	if err != nil {
		return errors.Wrap(err, "save config")
	}

	return nil
}
