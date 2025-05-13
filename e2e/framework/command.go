package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"

	"dev.khulnasoft.com/pkg/client"
	provider2 "dev.khulnasoft.com/pkg/provider"
	"dev.khulnasoft.com/pkg/workspace"
)

func (f *Framework) FindWorkspace(ctx context.Context, id string) (*provider2.Workspace, error) {
	list, err := f.DevSpaceListParsed(ctx)
	if err != nil {
		return nil, err
	}

	workspaceID := workspace.ToID(id)
	for _, w := range list {
		if w.ID == workspaceID {
			return w, nil
		}
	}

	return nil, fmt.Errorf("couldn't find workspace %s", workspaceID)
}

func (f *Framework) DevSpaceListParsed(ctx context.Context) ([]*provider2.Workspace, error) {
	raw, err := f.DevSpaceList(ctx)
	if err != nil {
		return nil, err
	}

	retList := []*provider2.Workspace{}
	err = json.Unmarshal([]byte(raw), &retList)
	if err != nil {
		return nil, err
	}

	return retList, nil
}

// DevSpaceList executes the `devspace list` command in the test framework
func (f *Framework) DevSpaceList(ctx context.Context) (string, error) {
	listArgs := []string{"list", "--output", "json"}

	out, _, err := f.ExecCommandCapture(ctx, listArgs)
	if err != nil {
		return "", fmt.Errorf("devspace list failed: %s", err.Error())
	}
	return out, nil
}

func (f *Framework) DevSpaceUpStreams(ctx context.Context, workspace string, additionalArgs ...string) (string, string, error) {
	upArgs := []string{"up", "--ide", "none", workspace}
	upArgs = append(upArgs, additionalArgs...)

	stdout, stderr, err := f.ExecCommandCapture(ctx, upArgs)
	if err != nil {
		return stdout, stderr, fmt.Errorf("devspace up failed: %s", err.Error())
	}

	return stdout, stderr, nil
}

// DevSpaceUp executes the `devspace up` command in the test framework
func (f *Framework) DevSpaceUpWithIDE(ctx context.Context, additionalArgs ...string) error {
	upArgs := []string{"up", "--debug"}
	upArgs = append(upArgs, additionalArgs...)

	_, _, err := f.ExecCommandCapture(ctx, upArgs)
	if err != nil {
		return fmt.Errorf("devspace up failed: %s", err.Error())
	}
	return nil
}

func (f *Framework) DevSpaceBuild(ctx context.Context, additionalArgs ...string) error {
	upArgs := []string{"build", "--debug"}
	upArgs = append(upArgs, additionalArgs...)

	_, _, err := f.ExecCommandCapture(ctx, upArgs)
	if err != nil {
		return fmt.Errorf("devspace build failed: %s", err.Error())
	}
	return nil
}

func (f *Framework) DevSpaceUp(ctx context.Context, additionalArgs ...string) error {
	upArgs := []string{"up", "--debug", "--ide", "none"}
	upArgs = append(upArgs, additionalArgs...)

	_, _, err := f.ExecCommandCapture(ctx, upArgs)
	if err != nil {
		return fmt.Errorf("devspace up failed: %s", err.Error())
	}
	return nil
}

func (f *Framework) DevSpaceUpRecreate(ctx context.Context, additionalArgs ...string) error {
	upArgs := []string{"up", "--recreate", "--debug", "--ide", "none"}
	upArgs = append(upArgs, additionalArgs...)

	_, _, err := f.ExecCommandCapture(ctx, upArgs)
	if err != nil {
		return fmt.Errorf("devspace up --recreate failed: %s", err.Error())
	}
	return nil
}

func (f *Framework) DevSpaceUpReset(ctx context.Context, additionalArgs ...string) error {
	upArgs := []string{"up", "--reset", "--debug", "--ide", "none"}
	upArgs = append(upArgs, additionalArgs...)

	_, _, err := f.ExecCommandCapture(ctx, upArgs)
	if err != nil {
		return fmt.Errorf("devspace up --reset failed: %s", err.Error())
	}
	return nil
}

func (f *Framework) DevSpaceSSH(ctx context.Context, workspace string, command string) (string, error) {
	out, err := f.ExecCommandOutput(ctx, []string{"ssh", workspace, "--command", command})
	if err != nil {
		return "", fmt.Errorf("devspace ssh failed: %s", err.Error())
	}
	return out, nil
}

func (f *Framework) DevSpaceSSHEchoTestString(ctx context.Context, workspace string) error {
	err := f.ExecCommand(ctx, true, true, "mYtEsTsTrInG", []string{"ssh", "--command", "echo 'bVl0RXNUc1RySW5H' | base64 -d", workspace})
	if err != nil {
		return fmt.Errorf("devspace ssh failed: %s", err.Error())
	}
	return nil
}

func (f *Framework) DevSpaceProviderOptionsCheckNamespaceDescription(ctx context.Context, provider, searchStr string) error {
	err := f.ExecCommand(ctx, true, true, searchStr, []string{"provider", "options", provider})
	if err != nil {
		return fmt.Errorf("did not found value %s in devspace provider options output. error: %s", searchStr, err.Error())
	}
	return nil
}

func (f *Framework) DevSpaceProviderList(ctx context.Context, extraArgs ...string) error {
	baseArgs := []string{"provider", "list"}
	err := f.ExecCommand(ctx, false, true, "", append(baseArgs, extraArgs...))
	if err != nil {
		return fmt.Errorf("devspace provider list failed: %s", err.Error())
	}
	return nil
}

func (f *Framework) DevSpaceProviderUse(ctx context.Context, provider string, extraArgs ...string) error {
	baseArgs := []string{"provider", "use", provider}
	err := f.ExecCommand(ctx, false, true, "", append(baseArgs, extraArgs...))
	if err != nil {
		return fmt.Errorf("devspace provider use failed: %s", err.Error())
	}
	return nil
}

func (f *Framework) DevSpaceStatus(ctx context.Context, extraArgs ...string) (client.WorkspaceStatus, error) {
	baseArgs := []string{"status", "--output", "json"}
	baseArgs = append(baseArgs, extraArgs...)
	stdout, err := f.ExecCommandOutput(ctx, baseArgs)
	if err != nil {
		return client.WorkspaceStatus{}, fmt.Errorf("devspace status failed: %s", err.Error())
	}

	status := &client.WorkspaceStatus{}
	err = json.Unmarshal([]byte(stdout), status)
	if err != nil {
		return client.WorkspaceStatus{}, err
	}

	return *status, nil
}

func (f *Framework) DevSpaceStop(ctx context.Context, workspace string) error {
	baseArgs := []string{"stop"}
	baseArgs = append(baseArgs, workspace)
	err := f.ExecCommand(ctx, false, false, "", baseArgs)
	if err != nil {
		return fmt.Errorf("devspace stop failed: %s", err.Error())
	}
	return nil
}

func (f *Framework) DevSpaceProviderAdd(ctx context.Context, args ...string) error {
	baseArgs := []string{"provider", "add"}
	baseArgs = append(baseArgs, args...)
	err := f.ExecCommand(ctx, false, false, "", baseArgs)
	if err != nil {
		return fmt.Errorf("devspace provider add failed: %s", err.Error())
	}
	return nil
}

func (f *Framework) DevSpaceProviderDelete(ctx context.Context, args ...string) error {
	baseArgs := []string{"provider", "delete"}
	baseArgs = append(baseArgs, args...)
	err := f.ExecCommand(ctx, false, false, "", baseArgs)
	if err != nil {
		return err
	}

	return nil
}

func (f *Framework) DevSpaceProviderUpdate(ctx context.Context, args ...string) error {
	baseArgs := []string{"provider", "update"}
	baseArgs = append(baseArgs, args...)
	err := f.ExecCommand(ctx, false, false, "", baseArgs)
	if err != nil {
		return fmt.Errorf("devspace provider update failed: %s", err.Error())
	}
	return nil
}

func (f *Framework) DevSpaceMachineCreate(args []string) error {
	baseArgs := []string{"machine", "create"}
	baseArgs = append(baseArgs, args...)
	err := f.ExecCommand(context.Background(), false, false, "", baseArgs)
	if err != nil {
		return fmt.Errorf("devspace nachine create failed: %s", err.Error())
	}
	return nil
}

func (f *Framework) DevSpaceMachineDelete(args []string) error {
	baseArgs := []string{"machine", "delete"}
	baseArgs = append(baseArgs, args...)
	err := f.ExecCommand(context.Background(), false, false, "", baseArgs)
	if err != nil {
		return fmt.Errorf("devspace nachine delete failed: %s", err.Error())
	}
	return nil
}

func (f *Framework) DevSpaceWorkspaceStop(ctx context.Context, extraArgs ...string) error {
	baseArgs := []string{"stop"}
	baseArgs = append(baseArgs, extraArgs...)
	return f.ExecCommandStdout(ctx, baseArgs)
}

func (f *Framework) DevSpaceWorkspaceDelete(ctx context.Context, workspace string, extraArgs ...string) error {
	baseArgs := []string{"delete", workspace, "--ignore-not-found"}
	baseArgs = append(baseArgs, extraArgs...)

	return f.ExecCommand(ctx, false, true, fmt.Sprintf("Successfully deleted workspace '%s'", workspace), baseArgs)
}

func (f *Framework) SetupGPG(tmpDir string) error {
	if _, err := exec.LookPath("gpg"); err != nil {
		err := exec.Command("sudo", "apt-get", " install", "gnupg2", "-y").Run()
		if err != nil {
			return nil
		}
	}

	err := exec.Command("gpg", "--import", filepath.Join(tmpDir, "gpg-public.key")).Run()
	if err != nil {
		return nil
	}

	err = exec.Command("gpg", "--import", filepath.Join(tmpDir, "gpg-private.key")).Run()
	if err != nil {
		return nil
	}

	err = exec.Command("gpgconf", "--kill", "gpg-agent").Run()
	if err != nil {
		return nil
	}

	err = exec.Command("gpg-agent", "--homedir", "$HOME/.gnupg", "--use-standard-socket", "--daemon").Run()
	if err != nil {
		return nil
	}

	return exec.Command("gpg", "-k").Run()
}

func (f *Framework) DevSpaceSSHGpgTestKey(ctx context.Context, workspace string) error {
	pubKeyB, err := exec.Command("sh", "-c", "gpg -k --with-colons 2>/dev/null | grep sec | base64 -w0").Output()
	if err != nil {
		return err
	}

	// First run to trigger the first forwarding
	stdout, _, err := f.ExecCommandCapture(ctx, []string{
		"ssh",
		"--agent-forwarding",
		"--gpg-agent-forwarding",
		"--command",
		"gpg -k --with-colons 2>/dev/null |grep sec |  base64 -w0", workspace,
	})
	if err != nil {
		return err
	}

	if stdout != string(pubKeyB) {
		return fmt.Errorf("devspace gpg public key forwarding failed, expected %s, got %s", string(pubKeyB), stdout)
	}

	return nil
}

func (f *Framework) DevspacePortTest(ctx context.Context, port string, workspace string) error {
	// First run to trigger the first forwarding
	_, _, err := f.ExecCommandCapture(ctx, []string{
		"ssh",
		"--forward-ports", port, workspace,
	})
	return err
}

func (f *Framework) DevSpaceProviderFindOption(ctx context.Context, provider string, searchStr string, extraArgs ...string) error {
	baseArgs := []string{"provider", "options", provider}
	err := f.ExecCommand(ctx, false, true, searchStr, append(baseArgs, extraArgs...))
	if err != nil {
		return fmt.Errorf("devspace provider use failed: %s", err.Error())
	}
	return nil
}

func (f *Framework) DevSpaceContextCreate(ctx context.Context, name string, extraArgs ...string) error {
	baseArgs := []string{"context", "create", name}
	err := f.ExecCommand(ctx, false, true, "", append(baseArgs, extraArgs...))
	if err != nil {
		return fmt.Errorf("devspace context create failed: %s", err.Error())
	}
	return nil
}

func (f *Framework) DevSpaceContextUse(ctx context.Context, name string, extraArgs ...string) error {
	baseArgs := []string{"context", "use", name}
	err := f.ExecCommand(ctx, false, true, "", append(baseArgs, extraArgs...))
	if err != nil {
		return fmt.Errorf("devspace context use failed: %s", err.Error())
	}
	return nil
}

func (f *Framework) DevSpaceContextDelete(ctx context.Context, name string, extraArgs ...string) error {
	baseArgs := []string{"context", "delete", name}
	err := f.ExecCommand(ctx, false, true, "", append(baseArgs, extraArgs...))
	if err != nil {
		return fmt.Errorf("devspace context delete failed: %s", err.Error())
	}
	return nil
}
