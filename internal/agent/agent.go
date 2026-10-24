package agent

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/SudiMango/podlet/internal/agent/client"
	"github.com/SudiMango/podlet/internal/agent/config"
	"github.com/SudiMango/podlet/internal/contract"
)

type Agent struct {
	config *config.Config
	client *client.Client
}

func New() (*Agent, error) {
	agentConfig, err := config.LoadAgentConfig()
	if err != nil {
		return nil, err
	}

	return &Agent{
		config: agentConfig,
		client: client.New(agentConfig.APIURL, agentConfig.Token),
	}, nil
}

func (a *Agent) Heartbeat(ctx context.Context) (*contract.HeartbeatResponse, error) {
	var response contract.HeartbeatResponse
	if err := a.client.DoJSON(
		ctx,
		http.MethodPost,
		"/agent/heartbeat",
		nil,
		&response,
	); err != nil {
		return nil, err
	}

	return &response, nil
}

func (a *Agent) reportHeartbeat(ctx context.Context) (bool, error) {
	response, err := a.Heartbeat(ctx)
	if err != nil {
		return false, fmt.Errorf("heartbeat: %w", err)
	}
	if response.Action == nil {
		return false, nil
	}

	switch *response.Action {
	case contract.ActionUninstall:
		if err := a.uninstall(ctx); err != nil {
			return false, err
		}
		return true, nil
	default:
		return false, fmt.Errorf("unsupported heartbeat action %q", *response.Action)
	}
}

func (a *Agent) Run(ctx context.Context) error {
	uninstalled, err := a.reportHeartbeat(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	} else if uninstalled {
		return nil
	}

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			uninstalled, err := a.reportHeartbeat(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}
			if uninstalled {
				return nil
			}
		}
	}
}

func (a *Agent) uninstall(ctx context.Context) error {
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("find home directory: %w", err)
		}
		dataDir = filepath.Join(homeDir, ".local", "share")
	}

	uninstallerPath := filepath.Join(dataDir, "podlet", "uninstall-server-agent.sh")
	info, err := os.Stat(uninstallerPath)
	if err != nil {
		return fmt.Errorf("access agent uninstaller: %w", err)
	}
	if info.Mode().Perm()&0111 == 0 {
		return fmt.Errorf("agent uninstaller is not executable: %s", uninstallerPath)
	}

	if err := a.client.DoJSON(ctx, http.MethodDelete, "/agent", nil, nil); err != nil {
		return fmt.Errorf("delete server registration: %w", err)
	}

	uninstaller := exec.CommandContext(ctx, uninstallerPath, "--from-agent")
	uninstaller.Stdout = os.Stdout
	uninstaller.Stderr = os.Stderr
	if err := uninstaller.Run(); err != nil {
		return fmt.Errorf("run agent uninstaller: %w", err)
	}

	return nil
}
