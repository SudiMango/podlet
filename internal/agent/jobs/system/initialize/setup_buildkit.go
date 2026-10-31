package initialize

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/SudiMango/podlet/internal/contract"
)

type SetupBuildkitJob struct{}

func (SetupBuildkitJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	output, err := exec.CommandContext(
		ctx,
		"docker", "run", "-d",
		"--name", "podlet-buildkit",
		"--privileged",
		"--restart", "unless-stopped",
		"--label", "dev.podlet.managed=true",
		"--label", "dev.podlet.component=buildkit",
		"moby/buildkit:latest",
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("set up BuildKit: %w: %s", err, strings.TrimSpace(string(output)))
	}

	if err := os.Setenv("BUILDKIT_HOST", buildkitHost); err != nil {
		return fmt.Errorf("set BUILDKIT_HOST: %w", err)
	}

	return nil
}
