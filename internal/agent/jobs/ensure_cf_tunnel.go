package jobs

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/agent/client"
	"github.com/SudiMango/podlet/internal/contract"
)

var (
	ErrCfTunnelNotFound = errors.New("Tunnel not found.")
)

type EnsureCfTunnelJob struct {
	client *client.Client
}

func NewEnsureCfTunnelJob(client *client.Client) EnsureCfTunnelJob {
	return EnsureCfTunnelJob{
		client: client,
	}
}

func (j EnsureCfTunnelJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	fmt.Println("Checking for Cloudflare tunnel...")

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

	switch {
	case err == nil:
		fmt.Println("Cloudflare tunnel found!")
		return nil

	case client.IsStatus(err, http.StatusNotFound):
		return ErrCfTunnelNotFound

	default:
		return fmt.Errorf("check for Cloudflare tunnel: %w", err)
	}
}
