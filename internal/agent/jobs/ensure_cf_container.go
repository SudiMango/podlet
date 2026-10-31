package jobs

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/SudiMango/podlet/internal/contract"
)

type EnsureCfContainerJob struct{}

func (EnsureCfContainerJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	output, err := exec.CommandContext(ctx, "docker", "container", "inspect", "podlet-cloudflared").CombinedOutput()
	if err != nil {
		return fmt.Errorf("check for Cloudflare container: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return nil
}
