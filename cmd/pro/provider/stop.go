package provider

import (
	"context"
	"fmt"
	"io"
	"os"

	storagev1 "dev.khulnasoft.com/api/v4/pkg/apis/storage/v1"

	"dev.khulnasoft.com/cmd/pro/flags"
	"dev.khulnasoft.com/pkg/platform"
	"dev.khulnasoft.com/pkg/platform/client"
	"dev.khulnasoft.com/pkg/platform/remotecommand"
	"dev.khulnasoft.com/log"
	"github.com/spf13/cobra"
)

// StopCmd holds the cmd flags
type StopCmd struct {
	*flags.GlobalFlags

	Log log.Logger
}

// NewStopCmd creates a new command
func NewStopCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &StopCmd{
		GlobalFlags: globalFlags,
		Log:         log.GetInstance(),
	}
	c := &cobra.Command{
		Hidden: true,
		Use:    "stop",
		Short:  "Runs stop on a workspace",
		Args:   cobra.NoArgs,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.Run(cobraCmd.Context(), os.Stdin, os.Stdout, os.Stderr)
		},
	}

	return c
}

func (cmd *StopCmd) Run(ctx context.Context, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	baseClient, err := client.InitClientFromPath(ctx, cmd.Config)
	if err != nil {
		return err
	}

	info, err := platform.GetWorkspaceInfoFromEnv()
	if err != nil {
		return err
	}
	workspace, err := platform.FindInstanceInProject(ctx, baseClient, info.UID, info.ProjectName)
	if err != nil {
		return err
	} else if workspace == nil {
		return fmt.Errorf("couldn't find workspace")
	}

	conn, err := platform.DialInstance(baseClient, workspace, "stop", platform.OptionsFromEnv(storagev1.DevPodFlagsStop), cmd.Log)
	if err != nil {
		return err
	}

	_, err = remotecommand.ExecuteConn(ctx, conn, stdin, stdout, stderr, cmd.Log.ErrorStreamOnly())
	if err != nil {
		return fmt.Errorf("error executing: %w", err)
	}

	return nil
}
