package jobs

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/SudiMango/podlet/internal/agent/config"
	"github.com/SudiMango/podlet/internal/contract"
)

type CreateCfContainerJob struct {
	config *config.Config
}

func NewCreateCfContainerJob(config *config.Config) CreateCfContainerJob {
	return CreateCfContainerJob{config: config}
}

func (j CreateCfContainerJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	fmt.Println("Creating Cloudflare container...")

	if j.config == nil || j.config.CfConnectorToken == "" {
		return errors.New("create Cloudflare container: connector token is missing")
	}

	output, err := exec.CommandContext(
		ctx,
		"docker", "run", "-d",
		"--name", "podlet-cloudflared",
		"--network", "podlet",
		"--restart", "unless-stopped",
		"--label", "dev.podlet.managed=true",
		"--label", "dev.podlet.component=cloudflared",
		"cloudflare/cloudflared:latest",
		"tunnel", "--no-autoupdate", "run",
		"--token", j.config.CfConnectorToken,
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("create Cloudflare container: %w: %s", err, strings.TrimSpace(string(output)))
	}

	fmt.Println("Successfully created Cloudflare container!")
	return nil
}
