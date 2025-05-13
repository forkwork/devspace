package ide

import (
	"dev.khulnasoft.com/cmd/flags"
	"github.com/spf13/cobra"
)

// NewIDECmd returns a new command
func NewIDECmd(flags *flags.GlobalFlags) *cobra.Command {
	ideCmd := &cobra.Command{
		Use:   "ide",
		Short: "DevSpace IDE commands",
	}

	ideCmd.AddCommand(NewUseCmd(flags))
	ideCmd.AddCommand(NewSetOptionsCmd(flags))
	ideCmd.AddCommand(NewOptionsCmd(flags))
	ideCmd.AddCommand(NewListCmd(flags))
	return ideCmd
}
