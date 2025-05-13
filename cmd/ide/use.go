package ide

import (
	"context"
	"fmt"
	"strings"

	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/ide"
	"dev.khulnasoft.com/pkg/ide/ideparse"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// UseCmd holds the use cmd flags
type UseCmd struct {
	flags.GlobalFlags

	Options []string
}

// NewUseCmd creates a new command
func NewUseCmd(flags *flags.GlobalFlags) *cobra.Command {
	cmd := &UseCmd{
		GlobalFlags: *flags,
	}
	useCmd := &cobra.Command{
		Use:   "use",
		Short: "Configure the default IDE to use (list available IDEs with 'devspace ide list')",
		Long: `Configure the default IDE to use

Available IDEs can be listed with 'devspace ide list'`,
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("please specify the ide to use, list available IDEs with 'devspace ide list'")
			}

			return cmd.Run(context.Background(), args[0])
		},
	}

	useCmd.Flags().StringArrayVarP(&cmd.Options, "option", "o", []string{}, "IDE option in the form KEY=VALUE")
	return useCmd
}

// Run runs the command logic
func (cmd *UseCmd) Run(ctx context.Context, ide string) error {
	devSpaceConfig, err := config.LoadConfig(cmd.Context, cmd.Provider)
	if err != nil {
		return err
	}

	ide = strings.ToLower(ide)
	ideOptions, err := ideparse.GetIDEOptions(ide)
	if err != nil {
		return err
	}

	// check if there are user options set
	if len(cmd.Options) > 0 {
		err = setOptions(devSpaceConfig, ide, cmd.Options, ideOptions)
		if err != nil {
			return err
		}
	}

	devSpaceConfig.Current().DefaultIDE = ide
	err = config.SaveConfig(devSpaceConfig)
	if err != nil {
		return errors.Wrap(err, "save config")
	}

	return nil
}

func setOptions(devSpaceConfig *config.Config, ide string, options []string, ideOptions ide.Options) error {
	optionValues, err := ideparse.ParseOptions(options, ideOptions)
	if err != nil {
		return err
	}

	if devSpaceConfig.Current().IDEs == nil {
		devSpaceConfig.Current().IDEs = map[string]*config.IDEConfig{}
	}

	newValues := map[string]config.OptionValue{}
	if devSpaceConfig.Current().IDEs[ide] != nil {
		for k, v := range devSpaceConfig.Current().IDEs[ide].Options {
			newValues[k] = v
		}
	}
	for k, v := range optionValues {
		newValues[k] = v
	}

	devSpaceConfig.Current().IDEs[ide] = &config.IDEConfig{
		Options: newValues,
	}
	return nil
}
