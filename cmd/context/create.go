package context

import (
	"context"
	"fmt"
	"strings"

	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/pkg/config"
	provider2 "dev.khulnasoft.com/pkg/provider"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// CreateCmd holds the create cmd flags
type CreateCmd struct {
	flags.GlobalFlags

	Options []string
}

// NewCreateCmd creates a new command
func NewCreateCmd(flags *flags.GlobalFlags) *cobra.Command {
	cmd := &CreateCmd{
		GlobalFlags: *flags,
	}
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new DevSpace context",
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("please specify the context to create")
			}

			return cmd.Run(context.Background(), args[0])
		},
	}

	createCmd.Flags().StringArrayVarP(&cmd.Options, "option", "o", []string{}, "context option in the form KEY=VALUE")
	return createCmd
}

// Run runs the command logic
func (cmd *CreateCmd) Run(ctx context.Context, context string) error {
	devSpaceConfig, err := config.LoadConfig("", cmd.Provider)
	if err != nil {
		return err
	} else if devSpaceConfig.Contexts[context] != nil {
		return fmt.Errorf("context '%s' already exists", context)
	}

	// verify name
	if provider2.ProviderNameRegEx.MatchString(context) {
		return fmt.Errorf("context name can only include smaller case letters, numbers or dashes")
	} else if len(context) > 48 {
		return fmt.Errorf("context name cannot be longer than 48 characters")
	}
	devSpaceConfig.Contexts[context] = &config.ContextConfig{}

	// check if there are create options set
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

func setOptions(devSpaceConfig *config.Config, context string, options []string) error {
	optionValues, err := parseOptions(options)
	if err != nil {
		return err
	} else if devSpaceConfig.Contexts[context] == nil {
		return fmt.Errorf("context '%s' doesn't exist", context)
	}

	newValues := map[string]config.OptionValue{}
	if devSpaceConfig.Contexts[context].Options != nil {
		for k, v := range devSpaceConfig.Contexts[context].Options {
			newValues[k] = v
		}
	}
	for k, v := range optionValues {
		newValues[k] = v
	}

	devSpaceConfig.Contexts[context].Options = newValues
	return nil
}

func parseOptions(options []string) (map[string]config.OptionValue, error) {
	allowedOptions := []string{}
	contextOptions := map[string]config.ContextOption{}
	for _, option := range config.ContextOptions {
		allowedOptions = append(allowedOptions, option.Name)
		contextOptions[option.Name] = option
	}

	retMap := map[string]config.OptionValue{}
	for _, option := range options {
		splitted := strings.Split(option, "=")
		if len(splitted) == 1 {
			return nil, fmt.Errorf("invalid option '%s', expected format KEY=VALUE", option)
		}

		key := strings.ToUpper(strings.TrimSpace(splitted[0]))
		value := strings.Join(splitted[1:], "=")
		contextOption, ok := contextOptions[key]
		if !ok {
			return nil, fmt.Errorf("invalid option '%s', allowed options are: %v", key, allowedOptions)
		}

		if len(contextOption.Enum) > 0 {
			found := false
			for _, e := range contextOption.Enum {
				if value == e {
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("invalid value '%s' for option '%s', has to match one of the following values: %v", value, key, contextOption.Enum)
			}
		}

		retMap[key] = config.OptionValue{
			Value:        value,
			UserProvided: true,
		}
	}

	return retMap, nil
}
