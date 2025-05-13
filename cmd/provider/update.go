package provider

import (
	"context"
	"fmt"

	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/workspace"
	"dev.khulnasoft.com/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// UpdateCmd holds the cmd flags
type UpdateCmd struct {
	*flags.GlobalFlags

	Use     bool
	Options []string
}

// NewUpdateCmd creates a new command
func NewUpdateCmd(flags *flags.GlobalFlags) *cobra.Command {
	cmd := &UpdateCmd{
		GlobalFlags: flags,
	}
	updateCmd := &cobra.Command{
		Use:   "update [name] [URL or path]",
		Short: "Updates a provider in DevSpace",
		RunE: func(_ *cobra.Command, args []string) error {
			ctx := context.Background()
			devSpaceConfig, err := config.LoadConfig(cmd.Context, cmd.Provider)
			if err != nil {
				return err
			}

			return cmd.Run(ctx, devSpaceConfig, args)
		},
	}

	updateCmd.Flags().BoolVar(&cmd.Use, "use", true, "If enabled will automatically activate the provider")
	updateCmd.Flags().StringArrayVarP(&cmd.Options, "option", "o", []string{}, "Provider option in the form KEY=VALUE")
	return updateCmd
}

func (cmd *UpdateCmd) Run(ctx context.Context, devSpaceConfig *config.Config, args []string) error {
	if len(args) != 1 && len(args) != 2 {
		return fmt.Errorf("please specify either a local file, url or git repository. E.g. devspace provider update my-provider khulnasoftdevspace-provider-gcloud")
	}

	providerSource := ""
	if len(args) == 2 {
		providerSource = args[1]
	}

	providerConfig, err := workspace.UpdateProvider(devSpaceConfig, args[0], providerSource, log.Default)
	if err != nil {
		return err
	}

	log.Default.Donef("Successfully updated provider %s", providerConfig.Name)
	if cmd.Use {
		err = ConfigureProvider(ctx, providerConfig, devSpaceConfig.DefaultContext, cmd.Options, false, false, false, nil, log.Default)
		if err != nil {
			log.Default.Errorf("Error configuring provider, please retry with 'devspace provider use %s --reconfigure'", providerConfig.Name)
			return errors.Wrap(err, "configure provider")
		}

		return nil
	}

	log.Default.Infof("To use the provider, please run the following command:")
	log.Default.Infof("devspace provider use %s", providerConfig.Name)
	return nil
}
