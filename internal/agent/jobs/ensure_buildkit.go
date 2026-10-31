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
	fmt.Println("Checking for BuildKit container...")

	output, err := exec.CommandContext(ctx, "docker", "container", "inspect", "podlet-buildkit").CombinedOutput()
	if err != nil {
		return fmt.Errorf("check for BuildKit container: %w: %s", err, strings.TrimSpace(string(output)))
	}

	if err := os.Setenv("BUILDKIT_HOST", buildkitHost); err != nil {
		return fmt.Errorf("set BUILDKIT_HOST: %w", err)
	}

	fmt.Println("BuildKit container found!")
	return nil
}
