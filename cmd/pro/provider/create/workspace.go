package create

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	managementv1 "dev.khulnasoft.com/api/v4/pkg/apis/management/v1"
	"dev.khulnasoft.com/cmd/pro/flags"
	"dev.khulnasoft.com/pkg/platform"
	"dev.khulnasoft.com/pkg/platform/client"
	"dev.khulnasoft.com/pkg/platform/form"
	"dev.khulnasoft.com/pkg/platform/project"
	"dev.khulnasoft.com/pkg/provider"
	"dev.khulnasoft.com/log"
	"dev.khulnasoft.com/log/terminal"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WorkspaceCmd holds the cmd flags
type WorkspaceCmd struct {
	*flags.GlobalFlags

	Log log.Logger
}

// NewWorkspaceCmd creates a new command
func NewWorkspaceCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &WorkspaceCmd{
		GlobalFlags: globalFlags,
		Log:         log.GetInstance().ErrorStreamOnly(),
	}
	c := &cobra.Command{
		Use:    "workspace",
		Short:  "Create a workspace",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.Run(cobraCmd.Context(), os.Stdin, os.Stdout, os.Stderr)
		},
	}

	return c
}

func (cmd *WorkspaceCmd) Run(ctx context.Context, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	baseClient, err := client.InitClientFromPath(ctx, cmd.Config)
	if err != nil {
		return err
	}

	// fully serialized intance, right now only used by GUI
	instanceEnv := os.Getenv(platform.WorkspaceInstanceEnv)
	if instanceEnv != "" {
		instance := &managementv1.DevSpaceWorkspaceInstance{} // init pointer
		err := json.Unmarshal([]byte(instanceEnv), instance)
		if err != nil {
			return fmt.Errorf("unmarshal workpace instance %s: %w", instanceEnv, err)
		}

		updatedInstance, err := createInstance(ctx, baseClient, instance, cmd.Log)
		if err != nil {
			return err
		}

		out, err := json.Marshal(updatedInstance)
		if err != nil {
			return err
		}

		fmt.Println(string(out))
		return nil
	}

	// Info through env, right now only used by CLI
	workspaceID := os.Getenv(provider.WORKSPACE_ID)
	workspaceUID := os.Getenv(provider.WORKSPACE_UID)
	workspaceFolder := os.Getenv(provider.WORKSPACE_FOLDER)
	workspaceContext := os.Getenv(provider.WORKSPACE_CONTEXT)
	workspacePicture := os.Getenv(platform.WorkspacePictureEnv)
	workspaceSource := os.Getenv(platform.WorkspaceSourceEnv)
	if workspaceUID == "" || workspaceID == "" || workspaceFolder == "" {
		return fmt.Errorf("workspaceID, workspaceUID or workspace folder not found: %s, %s, %s", workspaceID, workspaceUID, workspaceFolder)
	}
	instance, err := platform.FindInstance(ctx, baseClient, workspaceUID)
	if err != nil {
		return err
	}
	// Nothing left to do if we already have an instance
	if instance != nil {
		return nil
	}
	if !terminal.IsTerminalIn {
		return fmt.Errorf("unable to create new instance through CLI if stdin is not a terminal")
	}

	instance, err = form.CreateInstance(ctx, baseClient, workspaceID, workspaceUID, workspaceSource, workspacePicture, cmd.Log)
	if err != nil {
		return err
	}

	_, err = createInstance(ctx, baseClient, instance, cmd.Log)
	if err != nil {
		return err
	}

	// once we have the instance, update workspace and save config
	// TODO: Do we need a file lock?
	workspaceConfig, err := provider.LoadWorkspaceConfig(workspaceContext, workspaceID)
	if err != nil {
		return fmt.Errorf("load workspace config: %w", err)
	}
	workspaceConfig.Pro = &provider.ProMetadata{
		InstanceName: instance.GetName(),
		Project:      project.ProjectFromNamespace(instance.GetNamespace()),
		DisplayName:  instance.Spec.DisplayName,
	}

	err = provider.SaveWorkspaceConfig(workspaceConfig)
	if err != nil {
		return fmt.Errorf("save workspace config: %w", err)
	}

	return nil
}

func createInstance(ctx context.Context, client client.Client, instance *managementv1.DevSpaceWorkspaceInstance, log log.Logger) (*managementv1.DevSpaceWorkspaceInstance, error) {
	managementClient, err := client.Management()
	if err != nil {
		return nil, err
	}

	updatedInstance, err := managementClient.Loft().ManagementV1().
		DevSpaceWorkspaceInstances(instance.GetNamespace()).
		Create(ctx, instance, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("create workspace instance: %w", err)
	}

	return platform.WaitForInstance(ctx, client, updatedInstance, log)
}
