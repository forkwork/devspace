package completion

import (
	"strings"

	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/platform"
	"dev.khulnasoft.com/pkg/workspace"
	"dev.khulnasoft.com/log"
	"github.com/spf13/cobra"
)

func GetPlatformHostSuggestions(rootCmd *cobra.Command, context, provider string, args []string, toComplete string, owner platform.OwnerFilter, logger log.Logger) ([]string, cobra.ShellCompDirective) {
	devSpaceConfig, err := config.LoadConfig(context, provider)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	proInstances, err := workspace.ListProInstances(devSpaceConfig, logger)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var suggestions []string

	for _, instance := range proInstances {
		if strings.HasPrefix(instance.Host, toComplete) {
			suggestions = append(suggestions, instance.Host)
		}
	}

	return suggestions, cobra.ShellCompDirectiveNoFileComp
}
