package update

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

	// GUI
	instanceEnv := os.Getenv(platform.WorkspaceInstanceEnv)
	if instanceEnv != "" {
		newInstance := &managementv1.DevPodWorkspaceInstance{}
		err := json.Unmarshal([]byte(instanceEnv), newInstance)
		if err != nil {
			return fmt.Errorf("unmarshal workpace instance %s: %w", instanceEnv, err)
		}
		newInstance.TypeMeta = metav1.TypeMeta{} // ignore

		projectName := project.ProjectFromNamespace(newInstance.GetNamespace())
		oldInstance, err := platform.FindInstanceByName(ctx, baseClient, newInstance.GetName(), projectName)
		if err != nil {
			return err
		}

		updatedInstance, err := updateInstance(ctx, baseClient, oldInstance, newInstance, cmd.Log)
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

	// CLI
	if !terminal.IsTerminalIn {
		return fmt.Errorf("unable to update instance through CLI if stdin is not a terminal")
	}
	workspaceID := os.Getenv(platform.WorkspaceIDEnv)
	workspaceUID := os.Getenv(platform.WorkspaceUIDEnv)
	project := os.Getenv(platform.ProjectEnv)
	if workspaceUID == "" || workspaceID == "" || project == "" {
		return fmt.Errorf("workspaceID, workspaceUID or project not found: %s, %s, %s", workspaceID, workspaceUID, project)
	}

	oldInstance, err := platform.FindInstanceInProject(ctx, baseClient, workspaceUID, project)
	if err != nil {
		return err
	}

	newInstance, err := form.UpdateInstance(ctx, baseClient, oldInstance, cmd.Log)
	if err != nil {
		return err
	}

	_, err = updateInstance(ctx, baseClient, oldInstance, newInstance, cmd.Log)
	if err != nil {
		return err
	}

	return nil
}

func updateInstance(ctx context.Context, client client.Client, oldInstance *managementv1.DevPodWorkspaceInstance, newInstance *managementv1.DevPodWorkspaceInstance, log log.Logger) (*managementv1.DevPodWorkspaceInstance, error) {
	// This ensures the template is kept up to date with configuration changes
	if newInstance.Spec.TemplateRef != nil {
		newInstance.Spec.TemplateRef.SyncOnce = true
	}

	return platform.UpdateInstance(ctx, client, oldInstance, newInstance, log)
}
