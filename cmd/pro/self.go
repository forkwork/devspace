package pro

import (
	"bytes"
	"context"
	"fmt"

	"dev.khulnasoft.com/cmd/pro/flags"
	"dev.khulnasoft.com/pkg/client/clientimplementation"
	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/provider"
	"dev.khulnasoft.com/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// SelfCmd holds the cmd flags
type SelfCmd struct {
	*flags.GlobalFlags
	Log log.Logger

	Host string
}

// NewSelfCmd creates a new command
func NewSelfCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &SelfCmd{
		GlobalFlags: globalFlags,
		Log:         log.GetInstance(),
	}
	c := &cobra.Command{
		Use:    "self",
		Short:  "Get self",
		Hidden: true,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			devSpaceConfig, provider, err := findProProvider(cobraCmd.Context(), cmd.Context, cmd.Provider, cmd.Host, cmd.Log)
			if err != nil {
				return err
			}

			return cmd.Run(cobraCmd.Context(), devSpaceConfig, provider)
		},
	}

	c.Flags().StringVar(&cmd.Host, "host", "", "The pro instance to use")
	_ = c.MarkFlagRequired("host")

	return c
}

func (cmd *SelfCmd) Run(ctx context.Context, devSpaceConfig *config.Config, provider *provider.ProviderConfig) error {
	var buf bytes.Buffer
	// ignore --debug because we tunnel json through stdio
	cmd.Log.SetLevel(logrus.InfoLevel)

	err := clientimplementation.RunCommandWithBinaries(
		ctx,
		"getSelf",
		provider.Exec.Proxy.Get.Self,
		devSpaceConfig.DefaultContext,
		nil,
		nil,
		devSpaceConfig.ProviderOptions(provider.Name),
		provider,
		nil,
		nil,
		&buf,
		nil,
		cmd.Log)
	if err != nil {
		return fmt.Errorf("get self: %w", err)
	}

	fmt.Println(buf.String())

	return nil
}
