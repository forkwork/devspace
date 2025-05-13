package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/blang/semver"
	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/provider"
	"dev.khulnasoft.com/pkg/workspace"
	"dev.khulnasoft.com/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	errProviderNotFound = errors.New("provider not found")
)

type CheckProviderUpdateCmd struct {
	*flags.GlobalFlags
	log log.Logger
}

type providerVersionCheck struct {
	UpdateAvailable bool   `json:"updateAvailable"`
	LatestVersion   string `json:"latestVersion,omitempty"`
}

// NewCheckProviderUpdateCmd creates a new command
func NewCheckProviderUpdateCmd(flags *flags.GlobalFlags) *cobra.Command {
	cmd := &CheckProviderUpdateCmd{
		GlobalFlags: flags,
		log:         log.Default,
	}
	shellCmd := &cobra.Command{
		Use:   "check-provider-update",
		Short: "Check if a provider update is available",
		RunE: func(_ *cobra.Command, args []string) error {
			devSpaceConfig, err := config.LoadConfig(cmd.Context, cmd.Provider)
			if err != nil {
				return err
			}
			return cmd.Run(context.Background(), devSpaceConfig, args)
		},
	}

	return shellCmd
}

func (cmd *CheckProviderUpdateCmd) Run(ctx context.Context, devSpaceConfig *config.Config, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("provider is missing")
	}
	providerName := args[0]

	providerSourceRaw, err := workspace.ResolveProviderSource(devSpaceConfig, providerName, cmd.log)
	if err != nil {
		return fmt.Errorf("provider %s doesn't exist", providerName)
	}

	// retrieve current config for provider
	allProviders, err := workspace.LoadAllProviders(devSpaceConfig, cmd.log)
	if err != nil {
		return err
	}
	currentProvider, ok := allProviders[providerName]
	if !ok {
		return errProviderNotFound
	}

	latestProviderConfig, err := loadLatestProvider(providerSourceRaw, cmd.log)
	if err != nil {
		return err
	}
	currentProviderVersion, err := semver.Parse(strings.TrimPrefix(currentProvider.Config.Version, "v"))
	if err != nil {
		return err
	}
	latestProviderVersion, err := semver.Parse(strings.TrimPrefix(latestProviderConfig.Version, "v"))
	if err != nil {
		return err
	}

	versionCheck := providerVersionCheck{UpdateAvailable: false}
	// check if new version is newer
	if latestProviderVersion.GT(currentProviderVersion) {
		versionCheck.UpdateAvailable = true
		versionCheck.LatestVersion = latestProviderConfig.Version
	}
	out, err := json.Marshal(versionCheck)
	if err != nil {
		return err
	}
	fmt.Println(string(out))

	return nil
}

func loadLatestProvider(providerSourceRaw string, log log.Logger) (*provider.ProviderConfig, error) {
	providerRaw, _, err := workspace.ResolveProvider(providerSourceRaw, log)
	if err != nil {
		return nil, errors.Wrap(err, "resolve provider")
	}
	providerConfig, err := provider.ParseProvider(bytes.NewReader(providerRaw))
	if err != nil {
		return nil, errors.Wrap(err, "parse provider")
	}

	return providerConfig, nil
}
