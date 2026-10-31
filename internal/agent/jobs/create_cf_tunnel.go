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

type CreateCfTunnelJob struct {
	client *client.Client
	config *config.Config
}

func NewCreateCfTunnelJob(client *client.Client, config *config.Config) CreateCfTunnelJob {
	return CreateCfTunnelJob{
		client: client,
		config: config,
	}
}

func (j CreateCfTunnelJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	if j.client == nil {
		return errors.New("create Cloudflare tunnel: client is not configured")
	}

	var tunnel contract.TunnelResponse
	if err := j.client.DoJSON(
		ctx,
		http.MethodPost,
		"/agent/cftunnel",
		nil,
		&tunnel,
	); err != nil {
		return fmt.Errorf("create Cloudflare tunnel: %w", err)
	}

	if tunnel.ConnectorToken == nil || *tunnel.ConnectorToken == "" {
		return errors.New("create Cloudflare tunnel: connector token is missing")
	}

	j.config.CfConnectorToken = *tunnel.ConnectorToken

	if err := config.SaveAgentConfig(j.config); err != nil {
		return fmt.Errorf("save connector token: %w", err)
	}

	return nil
}
