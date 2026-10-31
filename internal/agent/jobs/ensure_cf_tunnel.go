package jobs

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/agent/client"
	"github.com/SudiMango/podlet/internal/agent/config"
	"github.com/SudiMango/podlet/internal/contract"
)

var (
	ErrCfTunnelNotFound = errors.New("Tunnel not found.")
)

type EnsureCfTunnelJob struct {
	client *client.Client
	config *config.Config
}

func NewEnsureCfTunnelJob(client *client.Client, config *config.Config) EnsureCfTunnelJob {
	return EnsureCfTunnelJob{
		client: client,
		config: config,
	}
}

func (j EnsureCfTunnelJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	if j.client == nil {
		return errors.New("check for Cloudflare tunnel: client is not configured")
	}

	var tunnel contract.TunnelResponse
	err := j.client.DoJSON(
		ctx,
		http.MethodGet,
		"/agent/cftunnel",
		nil,
		&tunnel,
	)

	if err != nil {
		if client.IsStatus(err, http.StatusNotFound) {
			return ErrCfTunnelNotFound
		}
		return fmt.Errorf("check for Cloudflare tunnel: %w", err)
	}

	if tunnel.ConnectorToken == nil || *tunnel.ConnectorToken == "" {
		return errors.New("check for Cloudflare tunnel: connector token is missing")
	}
	if j.config == nil {
		return errors.New("check for Cloudflare tunnel: config is not configured")
	}

	j.config.CfConnectorToken = *tunnel.ConnectorToken
	if err := config.SaveAgentConfig(j.config); err != nil {
		return fmt.Errorf("save connector token: %w", err)
	}

	return nil
}
