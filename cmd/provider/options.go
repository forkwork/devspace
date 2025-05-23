package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"dev.khulnasoft.com/cmd/completion"
	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/types"
	"dev.khulnasoft.com/pkg/workspace"
	"dev.khulnasoft.com/log"
	"dev.khulnasoft.com/log/table"
	"github.com/spf13/cobra"
)

// OptionsCmd holds the options cmd flags
type OptionsCmd struct {
	*flags.GlobalFlags

	Hidden bool
	Output string
}

// NewOptionsCmd creates a new command
func NewOptionsCmd(flags *flags.GlobalFlags) *cobra.Command {
	cmd := &OptionsCmd{
		GlobalFlags: flags,
	}
	optionsCmd := &cobra.Command{
		Use:   "options [provider]",
		Short: "Show options of an existing provider",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.Run(context.Background(), args)
		},
		ValidArgsFunction: func(rootCmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completion.GetProviderSuggestions(rootCmd, cmd.Context, cmd.Provider, args, toComplete, cmd.Owner, log.Default)
		},
	}

	optionsCmd.Flags().BoolVar(&cmd.Hidden, "hidden", false, "If true, will also show hidden options.")
	optionsCmd.Flags().StringVar(&cmd.Output, "output", "plain", "The output format to use. Can be json or plain")
	return optionsCmd
}

type optionWithValue struct {
	types.Option `json:",inline"`

	Children []string `json:"children,omitempty"`
	Value    string   `json:"value,omitempty"`
}

// Run runs the command logic
func (cmd *OptionsCmd) Run(ctx context.Context, args []string) error {
	devSpaceConfig, err := config.LoadConfig(cmd.Context, cmd.Provider)
	if err != nil {
		return err
	}

	providerName := devSpaceConfig.Current().DefaultProvider
	if len(args) > 0 {
		providerName = args[0]
	} else if providerName == "" {
		return fmt.Errorf("please specify a provider")
	}

	if providerName != "" && cmd.GlobalFlags.Provider != "" {
		if providerName != cmd.GlobalFlags.Provider {
			log.Default.Infof("providerName=%+v", providerName)
			log.Default.Infof("GlobalFlags.Provider=%+v", cmd.GlobalFlags.Provider)
			return fmt.Errorf("ambiguous provider configuration detected")
		}
	}

	providerWithOptions, err := workspace.FindProvider(devSpaceConfig, providerName, log.Default.ErrorStreamOnly())
	if err != nil {
		return err
	}

	return printOptions(devSpaceConfig, providerWithOptions, cmd.Output, cmd.Hidden)
}

func printOptions(devSpaceConfig *config.Config, provider *workspace.ProviderWithOptions, format string, showHidden bool) error {
	entryOptions := devSpaceConfig.ProviderOptions(provider.Config.Name)
	dynamicOptions := devSpaceConfig.DynamicProviderOptionDefinitions(provider.Config.Name)
	srcOptions := MergeDynamicOptions(provider.Config.Options, dynamicOptions)
	if format == "plain" {
		tableEntries := [][]string{}
		for optionName, entry := range srcOptions {
			if !showHidden && entry.Hidden {
				continue
			}

			value := entryOptions[optionName].Value
			if value != "" && entry.Password {
				value = "********"
			}

			tableEntries = append(tableEntries, []string{
				optionName,
				strconv.FormatBool(entry.Required),
				entry.Description,
				entry.Default,
				value,
			})
		}
		sort.SliceStable(tableEntries, func(i, j int) bool {
			return tableEntries[i][0] < tableEntries[j][0]
		})

		table.PrintTable(log.Default, []string{
			"Name",
			"Required",
			"Description",
			"Default",
			"Value",
		}, tableEntries)
	} else if format == "json" {
		options := map[string]optionWithValue{}
		for optionName, entry := range srcOptions {
			if !showHidden && entry.Hidden {
				continue
			}

			options[optionName] = optionWithValue{
				Option:   *entry,
				Children: entryOptions[optionName].Children,
				Value:    entryOptions[optionName].Value,
			}
		}

		out, err := json.MarshalIndent(options, "", "  ")
		if err != nil {
			return err
		}
		fmt.Print(string(out))
	} else {
		return fmt.Errorf("unexpected output format, choose either json or plain. Got %s", format)
	}

	return nil
}

// MergeDynamicOptions merges the static provider options and dynamic options
func MergeDynamicOptions(options map[string]*types.Option, dynamicOptions config.OptionDefinitions) map[string]*types.Option {
	retOptions := map[string]*types.Option{}
	for k, option := range options {
		retOptions[k] = option
	}
	for k, option := range dynamicOptions {
		retOptions[k] = option
	}

	return retOptions
}
