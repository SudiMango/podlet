package jobs

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/SudiMango/podlet/internal/contract"
)

type EnsureDockerJob struct{}

func (EnsureDockerJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	fmt.Println("Checking for Docker...")

	output, err := exec.CommandContext(ctx, "docker", "--version").CombinedOutput()
	if err != nil {
		return fmt.Errorf("check for Docker: %w: %s", err, strings.TrimSpace(string(output)))
	}

	fmt.Println("Docker found!")
	return nil
}
