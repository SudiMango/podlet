package jobs

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/SudiMango/podlet/internal/contract"
)

const buildkitHost = "docker-container://podlet-buildkit"

type EnsureBuildkitJob struct{}

func (EnsureBuildkitJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	output, err := exec.CommandContext(
		ctx,
		"docker", "container", "inspect",
		"--format", "{{.State.Running}}",
		"podlet-buildkit",
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("check for BuildKit container: %w: %s", err, strings.TrimSpace(string(output)))
	}
	if strings.TrimSpace(string(output)) != "true" {
		output, err = exec.CommandContext(ctx, "docker", "start", "podlet-buildkit").CombinedOutput()
		if err != nil {
			return fmt.Errorf("start BuildKit container: %w: %s", err, strings.TrimSpace(string(output)))
		}
	}

	if err := os.Setenv("BUILDKIT_HOST", buildkitHost); err != nil {
		return fmt.Errorf("set BUILDKIT_HOST: %w", err)
	}

	return nil
}
