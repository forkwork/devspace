package cmd

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/extract"
	"dev.khulnasoft.com/pkg/provider"
	"dev.khulnasoft.com/pkg/workspace"
	"dev.khulnasoft.com/log"
	"github.com/spf13/cobra"
)

// ImportCmd holds the export cmd flags
type ImportCmd struct {
	*flags.GlobalFlags

	WorkspaceID string

	MachineID    string
	MachineReuse bool

	ProviderID    string
	ProviderReuse bool

	Data string
}

// NewImportCmd creates a new command
func NewImportCmd(flags *flags.GlobalFlags) *cobra.Command {
	cmd := &ImportCmd{
		GlobalFlags: flags,
	}
	importCmd := &cobra.Command{
		Use:    "import",
		Short:  "Imports a workspace configuration",
		Args:   cobra.NoArgs,
		Hidden: true,
		RunE: func(_ *cobra.Command, args []string) error {
			ctx := context.Background()
			devSpaceConfig, err := config.LoadConfig(cmd.Context, cmd.Provider)
			if err != nil {
				return err
			}

			return cmd.Run(ctx, devSpaceConfig, log.Default)
		},
	}

	importCmd.Flags().StringVar(&cmd.WorkspaceID, "workspace-id", "", "To workspace id to use")
	importCmd.Flags().StringVar(&cmd.MachineID, "machine-id", "", "The machine id to use")
	importCmd.Flags().BoolVar(&cmd.MachineReuse, "machine-reuse", false, "If machine already exists, reuse existing machine")
	importCmd.Flags().StringVar(&cmd.ProviderID, "provider-id", "", "The provider id to use")
	importCmd.Flags().BoolVar(&cmd.ProviderReuse, "provider-reuse", false, "If provider already exists, reuse existing provider")
	importCmd.Flags().StringVar(&cmd.Data, "data", "", "The data to import as raw json")
	_ = importCmd.MarkFlagRequired("data")
	return importCmd
}

// Run runs the command logic
func (cmd *ImportCmd) Run(ctx context.Context, devSpaceConfig *config.Config, log log.Logger) error {
	exportConfig := &provider.ExportConfig{}
	err := json.Unmarshal([]byte(cmd.Data), exportConfig)
	if err != nil {
		return fmt.Errorf("decode workspace data: %w", err)
	} else if exportConfig.Workspace == nil {
		return fmt.Errorf("workspace is missing in imported data")
	} else if exportConfig.Provider == nil {
		return fmt.Errorf("provider is missing in imported data")
	}

	// set ids correctly
	if cmd.MachineID == "" && exportConfig.Machine != nil {
		cmd.MachineID = exportConfig.Machine.ID
	}
	if cmd.WorkspaceID == "" {
		cmd.WorkspaceID = exportConfig.Workspace.ID
	}
	if cmd.ProviderID == "" {
		cmd.ProviderID = exportConfig.Provider.ID
	}

	// check if conflicting ids
	err = cmd.checkForConflictingIDs(ctx, exportConfig, devSpaceConfig, log)
	if err != nil {
		return err
	}

	// import provider
	err = cmd.importProvider(devSpaceConfig, exportConfig, log)
	if err != nil {
		return err
	}

	// import machine
	err = cmd.importMachine(devSpaceConfig, exportConfig, log)
	if err != nil {
		return err
	}

	// import workspace
	err = cmd.importWorkspace(devSpaceConfig, exportConfig, log)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *ImportCmd) importWorkspace(devSpaceConfig *config.Config, exportConfig *provider.ExportConfig, log log.Logger) error {
	workspaceDir, err := provider.GetWorkspaceDir(devSpaceConfig.DefaultContext, cmd.WorkspaceID)
	if err != nil {
		return fmt.Errorf("get workspace dir: %w", err)
	}

	err = os.MkdirAll(workspaceDir, 0755)
	if err != nil {
		return fmt.Errorf("create workspace dir: %w", err)
	}

	decoded, err := base64.RawStdEncoding.DecodeString(exportConfig.Workspace.Data)
	if err != nil {
		return fmt.Errorf("decode workspace data: %w", err)
	}

	err = extract.Extract(bytes.NewReader(decoded), workspaceDir)
	if err != nil {
		return fmt.Errorf("extract workspace data: %w", err)
	}

	// exchange config
	workspaceConfig, err := provider.LoadWorkspaceConfig(devSpaceConfig.DefaultContext, cmd.WorkspaceID)
	if err != nil {
		return fmt.Errorf("load machine config: %w", err)
	}
	workspaceConfig.ID = cmd.WorkspaceID
	workspaceConfig.Context = devSpaceConfig.DefaultContext
	workspaceConfig.Machine.ID = cmd.MachineID
	workspaceConfig.Provider.Name = cmd.ProviderID

	// save machine config
	err = provider.SaveWorkspaceConfig(workspaceConfig)
	if err != nil {
		return fmt.Errorf("save workspace config: %w", err)
	}

	log.Donef("Successfully imported workspace %s", cmd.WorkspaceID)
	return nil
}

func (cmd *ImportCmd) importMachine(devSpaceConfig *config.Config, exportConfig *provider.ExportConfig, log log.Logger) error {
	if exportConfig.Machine == nil {
		return nil
	}

	// if machine already exists we skip
	if cmd.MachineReuse && provider.MachineExists(devSpaceConfig.DefaultContext, cmd.MachineID) {
		log.Infof("Reusing existing machine %s", cmd.MachineID)
		return nil
	}

	machineDir, err := provider.GetMachineDir(devSpaceConfig.DefaultContext, cmd.MachineID)
	if err != nil {
		return fmt.Errorf("get machine dir: %w", err)
	}

	err = os.MkdirAll(machineDir, 0755)
	if err != nil {
		return fmt.Errorf("create machine dir: %w", err)
	}

	decoded, err := base64.RawStdEncoding.DecodeString(exportConfig.Machine.Data)
	if err != nil {
		return fmt.Errorf("decode machine data: %w", err)
	}

	err = extract.Extract(bytes.NewReader(decoded), machineDir)
	if err != nil {
		return fmt.Errorf("extract machine data: %w", err)
	}

	// exchange config
	machineConfig, err := provider.LoadMachineConfig(devSpaceConfig.DefaultContext, cmd.MachineID)
	if err != nil {
		return fmt.Errorf("load machine config: %w", err)
	}
	machineConfig.ID = cmd.MachineID
	machineConfig.Context = devSpaceConfig.DefaultContext
	machineConfig.Provider.Name = cmd.ProviderID

	// save machine config
	err = provider.SaveMachineConfig(machineConfig)
	if err != nil {
		return fmt.Errorf("save machine config: %w", err)
	}

	log.Donef("Successfully imported machine %s", cmd.MachineID)
	return nil
}

func (cmd *ImportCmd) importProvider(devSpaceConfig *config.Config, exportConfig *provider.ExportConfig, log log.Logger) error {
	// if provider already exists we skip
	if cmd.ProviderReuse && provider.ProviderExists(devSpaceConfig.DefaultContext, cmd.ProviderID) {
		log.Infof("Reusing existing provider %s", cmd.ProviderID)
		return nil
	}

	providerDir, err := provider.GetProviderDir(devSpaceConfig.DefaultContext, cmd.ProviderID)
	if err != nil {
		return fmt.Errorf("get provider dir: %w", err)
	}

	err = os.MkdirAll(providerDir, 0755)
	if err != nil {
		return fmt.Errorf("create provider dir: %w", err)
	}

	decoded, err := base64.RawStdEncoding.DecodeString(exportConfig.Provider.Data)
	if err != nil {
		return fmt.Errorf("decode provider data: %w", err)
	}

	err = extract.Extract(bytes.NewReader(decoded), providerDir)
	if err != nil {
		return fmt.Errorf("extract provider data: %w", err)
	}

	// exchange config
	providerConfig, err := provider.LoadProviderConfig(devSpaceConfig.DefaultContext, cmd.ProviderID)
	if err != nil {
		return fmt.Errorf("load provider config: %w", err)
	}
	providerConfig.Name = cmd.ProviderID

	// save provider config
	err = provider.SaveProviderConfig(devSpaceConfig.DefaultContext, providerConfig)
	if err != nil {
		return fmt.Errorf("save provider config: %w", err)
	}

	// add provider options
	if exportConfig.Provider.Config != nil {
		if devSpaceConfig.Current().Providers == nil {
			devSpaceConfig.Current().Providers = map[string]*config.ProviderConfig{}
		}

		devSpaceConfig.Current().Providers[cmd.ProviderID] = exportConfig.Provider.Config
		err = config.SaveConfig(devSpaceConfig)
		if err != nil {
			return fmt.Errorf("save devspace config: %w", err)
		}
	}

	log.Donef("Successfully imported provider %s", cmd.ProviderID)
	return nil
}

func (cmd *ImportCmd) checkForConflictingIDs(ctx context.Context, exportConfig *provider.ExportConfig, devSpaceConfig *config.Config, log log.Logger) error {
	workspaces, err := workspace.List(ctx, devSpaceConfig, false, cmd.Owner, log)
	if err != nil {
		return fmt.Errorf("error listing workspaces: %w", err)
	}

	// check for workspace duplicate
	if exportConfig.Workspace != nil {
		for _, workspace := range workspaces {
			if workspace.ID == cmd.WorkspaceID {
				return fmt.Errorf("existing workspace with id %s found, please use --workspace-id to override the workspace id", cmd.WorkspaceID)
			} else if workspace.UID == exportConfig.Workspace.UID {
				return fmt.Errorf("existing workspace %s with uid %s found, please use --workspace-id to override the workspace id", workspace.ID, workspace.UID)
			}
		}
	}

	// check if machine already exists
	if !cmd.MachineReuse && exportConfig.Machine != nil {
		if provider.MachineExists(devSpaceConfig.DefaultContext, cmd.MachineID) {
			return fmt.Errorf("existing machine with id %s found, please use --machine-reuse to skip importing the machine or --machine-id to override the machine id", cmd.MachineID)
		}
	}

	// check if provider already exists
	if !cmd.ProviderReuse && exportConfig.Provider != nil {
		if provider.ProviderExists(devSpaceConfig.DefaultContext, cmd.ProviderID) {
			return fmt.Errorf("existing provider with id %s found, please use --provider-reuse to skip importing the provider or --provider-id to override the provider id", cmd.ProviderID)
		}
	}

	return nil
}
