package provider

import (
	"os"

	"dev.khulnasoft.com/cmd/agent"
	"dev.khulnasoft.com/cmd/pro/flags"
	"dev.khulnasoft.com/cmd/pro/provider/create"
	"dev.khulnasoft.com/cmd/pro/provider/get"
	"dev.khulnasoft.com/cmd/pro/provider/list"
	"dev.khulnasoft.com/cmd/pro/provider/update"
	"dev.khulnasoft.com/cmd/pro/provider/watch"
	"dev.khulnasoft.com/pkg/client/clientimplementation"
	"dev.khulnasoft.com/pkg/platform"
	"dev.khulnasoft.com/pkg/platform/client"
	"dev.khulnasoft.com/pkg/telemetry"
	"dev.khulnasoft.com/log"

	"github.com/spf13/cobra"
)

// NewProProviderCmd creates a new cobra command
func NewProProviderCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	c := &cobra.Command{
		Use:    "provider",
		Short:  "DevPod Pro provider commands",
		Args:   cobra.NoArgs,
		Hidden: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if (globalFlags.Config == "" || globalFlags.Config == client.DefaultCacheConfig) && os.Getenv("LOFT_CONFIG") != "" {
				globalFlags.Config = os.Getenv(platform.ConfigEnv)
			}

			log.Default.SetFormat(log.JSONFormat)

			if os.Getenv(clientimplementation.DevPodDebug) == "true" {
				globalFlags.Debug = true
			}

			// Disable debug hints if we execute pro commands from DevPod Desktop
			// We're reusing the agent.AgentExecutedAnnotation for simplicity, could rename in the future
			if os.Getenv(telemetry.UIEnvVar) == "true" {
				cmd.VisitParents(func(c *cobra.Command) {
					// find the root command
					if c.Name() == "devspace" {
						if c.Annotations == nil {
							c.Annotations = map[string]string{}
						}
						c.Annotations[agent.AgentExecutedAnnotation] = "true"
					}
				})
			}
		},
	}

	c.AddCommand(list.NewCmd(globalFlags))
	c.AddCommand(watch.NewCmd(globalFlags))
	c.AddCommand(create.NewCmd(globalFlags))
	c.AddCommand(get.NewCmd(globalFlags))
	c.AddCommand(update.NewCmd(globalFlags))
	c.AddCommand(NewHealthCmd(globalFlags))

	c.AddCommand(NewUpCmd(globalFlags))
	c.AddCommand(NewStopCmd(globalFlags))
	c.AddCommand(NewSshCmd(globalFlags))
	c.AddCommand(NewStatusCmd(globalFlags))
	c.AddCommand(NewDeleteCmd(globalFlags))
	c.AddCommand(NewRebuildCmd(globalFlags))
	return c
}
