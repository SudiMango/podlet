package initialize

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/SudiMango/podlet/internal/contract"
)

type EnsureDockerJob struct{}

func (EnsureDockerJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	output, err := exec.CommandContext(ctx, "docker", "--version").CombinedOutput()
	if err != nil {
		return fmt.Errorf("check for Docker: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return nil
}
