package use

import (
	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/cmd/ide"
	"dev.khulnasoft.com/cmd/provider"
	"github.com/spf13/cobra"
)

// NewUseCmd returns a new root command
func NewUseCmd(flags *flags.GlobalFlags) *cobra.Command {
	useCmd := &cobra.Command{
		Use:   "use",
		Short: "Use DevSpace resources",
	}

	// use provider
	useProviderCmd := provider.NewUseCmd(flags)
	useProviderCmd.Use = "provider"
	useCmd.AddCommand(useProviderCmd)

	// use ide
	useIDECmd := ide.NewUseCmd(flags)
	useIDECmd.Use = "ide"
	useCmd.AddCommand(useIDECmd)
	return useCmd
}
