package jobs

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/SudiMango/podlet/internal/contract"
)

type InstallDockerJob struct{}

func (InstallDockerJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	fmt.Println("Installing Docker...")

	output, err := exec.CommandContext(
		ctx,
		"bash",
		"-o",
		"pipefail",
		"-c",
		"curl -fsSL https://get.docker.com | sh",
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("install Docker: %w: %s", err, strings.TrimSpace(string(output)))
	}

	fmt.Println("Successfully installed Docker!")
	return nil
}
