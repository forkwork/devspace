package ide

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/ide/ideparse"
	"dev.khulnasoft.com/log"
	"dev.khulnasoft.com/log/table"
	"github.com/spf13/cobra"
)

// ListCmd holds the list cmd flags
type ListCmd struct {
	flags.GlobalFlags

	Output string
}

// NewListCmd creates a new command
func NewListCmd(flags *flags.GlobalFlags) *cobra.Command {
	cmd := &ListCmd{
		GlobalFlags: *flags,
	}
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List available IDEs",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.Run(context.Background())
		},
	}

	listCmd.Flags().StringVar(&cmd.Output, "output", "plain", "The output format to use. Can be json or plain")
	return listCmd
}

type IDEWithDefault struct {
	ideparse.AllowedIDE `json:",inline"`

	Default bool `json:"default,omitempty"`
}

// Run runs the command logic
func (cmd *ListCmd) Run(ctx context.Context) error {
	devSpaceConfig, err := config.LoadConfig(cmd.Context, cmd.Provider)
	if err != nil {
		return err
	}

	if cmd.Output == "plain" {
		tableEntries := [][]string{}
		for _, entry := range ideparse.AllowedIDEs {
			tableEntries = append(tableEntries, []string{
				string(entry.Name),
				strconv.FormatBool(devSpaceConfig.Current().DefaultIDE == string(entry.Name)),
			})
		}
		sort.SliceStable(tableEntries, func(i, j int) bool {
			return tableEntries[i][0] < tableEntries[j][0]
		})

		table.PrintTable(log.Default, []string{
			"Name",
			"Default",
		}, tableEntries)
	} else if cmd.Output == "json" {
		ides := []IDEWithDefault{}
		for _, entry := range ideparse.AllowedIDEs {
			ides = append(ides, IDEWithDefault{
				AllowedIDE: entry,
				Default:    devSpaceConfig.Current().DefaultIDE == string(entry.Name),
			})
		}

		out, err := json.MarshalIndent(ides, "", "  ")
		if err != nil {
			return err
		}
		fmt.Print(string(out))
	} else {
		return fmt.Errorf("unexpected output format, choose either json or plain. Got %s", cmd.Output)
	}

	return nil
}
