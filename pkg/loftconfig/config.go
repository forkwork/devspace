package loftconfig

import (
	"fmt"
	"os/exec"

	"dev.khulnasoft.com/pkg/platform/client"
	"dev.khulnasoft.com/log"
)

func AuthDevspaceCliToPlatform(config *client.Config, logger log.Logger) error {
	cmd := exec.Command("devspace", "pro", "login", "--access-key", config.AccessKey, config.Host)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Debugf("Failed executing `devspace pro login`: %w, output: %s", err, out)
		return fmt.Errorf("error executing 'devspace pro login' command: %w, access key: %v, host: %v", err, config.AccessKey, config.Host)
	}

	return nil
}

func AuthVClusterCliToPlatform(config *client.Config, logger log.Logger) error {
	// Check if vcluster is available inside the workspace
	if _, err := exec.LookPath("vcluster"); err != nil {
		logger.Debugf("'vcluster' command is not available")
		return nil
	}

	cmd := exec.Command("vcluster", "login", "--access-key", config.AccessKey, config.Host)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Debugf("Failed executing `vcluster login` : %w, output: %s", err, out)
		return fmt.Errorf("error executing 'vcluster login' command: %w", err)
	}

	return nil
}
