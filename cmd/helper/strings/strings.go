package strings

import (
	"dev.khulnasoft.com/cmd/flags"
	"github.com/spf13/cobra"
)

// NewStringsCmd returns a new command
func NewStringsCmd(flags *flags.GlobalFlags) *cobra.Command {
	stringsCmd := &cobra.Command{
		Use:    "strings",
		Short:  "DevSpace String Utility Commands",
		Hidden: true,
	}

	return stringsCmd
}
