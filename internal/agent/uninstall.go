package agent

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/SudiMango/podlet/internal/agent/dto"
)

func (a *Agent) Uninstall(ctx context.Context) error {
	uninstallerPath, err := agentUninstallerPath()
	if err != nil {
		return err
	}

	info, err := os.Stat(uninstallerPath)
	if err != nil {
		return fmt.Errorf("access agent uninstaller: %w", err)
	}
	if info.Mode().Perm()&0111 == 0 {
		return fmt.Errorf("agent uninstaller is not executable: %s", uninstallerPath)
	}

	if err := a.removeManagedResources(ctx); err != nil {
		return fmt.Errorf("remove managed resources: %w", err)
	}

	if err := a.client.DoJSON(ctx, http.MethodDelete, "/agent", nil, nil); err != nil {
		return fmt.Errorf("delete server registration: %w", err)
	}

	cmd := exec.CommandContext(ctx, uninstallerPath, "--from-agent")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run agent uninstaller: %w", err)
	}

	return nil
}

func (a *Agent) removeManagedResources(ctx context.Context) error {
	jobTypes := []dto.JobType{
		dto.JobTypeRemoveCfContainer,
		dto.JobTypeRemoveBuildkit,
		dto.JobTypeRemovePodletNetwork,
		dto.JobTypeUninstallRailpack,
	}

	for _, jobType := range jobTypes {
		if err := a.executors[jobType].Execute(ctx, nil); err != nil {
			return fmt.Errorf("execute %s job: %w", jobType, err)
		}
	}

	return nil
}

func agentUninstallerPath() (string, error) {
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("find home directory: %w", err)
		}
		dataDir = filepath.Join(homeDir, ".local", "share")
	}

	return filepath.Join(dataDir, "podlet", "uninstall-agent.sh"), nil
}
